package client

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"math"
	"net"
	"sync"
	"time"

	cnf "github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func Client(config *cnf.AppConfigList) {
	// Address to send messages
	presenceNetworkEntrypointIp := config.Neighbors[0]
	pneIpString := fmt.Sprintf("%s:%d", presenceNetworkEntrypointIp.Ip, presenceNetworkEntrypointIp.Port)
	fmt.Println(pneIpString)

	listenIp := "0.0.0.0:2222"

	// Setup a UDP connection for sending
	conn, err := net.Dial("udp", pneIpString)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	// Create a UDP listener for incoming messages
	laddr := net.UDPAddr{
		Port: 2222,
		IP:   net.ParseIP(listenIp),
	}

	listener, err := net.ListenUDP("udp", &laddr)
	if err != nil {
		log.Fatalf("Failed to create UDP listener: %v", err)
	}

	// Create the app and window
	myApp := app.New()
	myWindow := myApp.NewWindow("Video Client")

	// Create an image canvas object
	imgCanvas := canvas.NewImageFromImage(nil)
	imgCanvas.FillMode = canvas.ImageFillContain

	// Create a label to display the msg.Path
	pathLabel := widget.NewLabel("")

	// Statistics variables
	var (
		prevDelay          float64
		jittersum          float64
		packetCount        int
		startTime          = time.Now()
		statsMutex         sync.Mutex
		totalSent          int
		totalBytesReceived int64
		frameDisplayCount  int
		fps                float64
		lastFpsUpdate      time.Time = time.Now()
		bufferMutex        sync.Mutex
		frameBuffer        []*image.Image
	)

	resetStats := func() {
		statsMutex.Lock()
		defer statsMutex.Unlock()

		// Reset statistics values
		packetCount = 0
		jittersum = 0.0
		prevDelay = 0.0
		totalBytesReceived = 0
		frameDisplayCount = 0
		fps = 0
		lastFpsUpdate = time.Now()
		startTime = time.Now()
	}

	// List of available videos
	videoOptions := []string{"lol", "demo", "video1.mp4", "video2.mp4"}

	// Entry widgets and buttons
	messageEntry := widget.NewEntry()
	messageEntry.SetPlaceHolder("Enter video, 'lol' or 'demo'")
	targetEntry := widget.NewEntry()
	targetEntry.SetPlaceHolder("Enter target, 's1'")

	// Dropdown widget for video selection
	videoSelect := widget.NewSelect(videoOptions, func(selected string) {
		messageEntry.SetText(selected) // Automatically set the selected video in the messageEntry
	})
	videoSelect.PlaceHolder = "Select a video"

	// Buttons
	sendButton := widget.NewButtonWithIcon("Send", theme.ConfirmIcon(), func() {
		sendCommand("PLAY", messageEntry.Text, targetEntry.Text, conn, config, listenIp, resetStats)
	})
	cancelButton := widget.NewButtonWithIcon("Stop", theme.CancelIcon(), func() {
		sendCommand("STOP", messageEntry.Text, targetEntry.Text, conn, config, listenIp, resetStats)
		resetStats()
	})
	statsButton := widget.NewButtonWithIcon("Statistics", theme.InfoIcon(), func() {
		showStatsWindow(myApp, startTime, &packetCount, &jittersum, &statsMutex, &totalSent, &totalBytesReceived, &fps)
	})

	// Layout
	controls := container.NewVBox(videoSelect, messageEntry, targetEntry, sendButton, cancelButton, statsButton, pathLabel)
	content := container.NewBorder(nil, controls, nil, nil, imgCanvas)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 600))

	// Goroutine for listening to incoming messages
	go func() {
		for {
			data, _, err := cnf.ReceiveMessageUDP(listener)

			if err != nil {
				log.Printf("Error receiving message: %s\n", err)
				continue
			}

			msg, err := cnf.UnmarshalHeader(data)
			if err != nil {
				log.Printf("Error unmarshaling message: %s\n", err)
				continue
			}

			if msg.GetServerVideoChunk() == nil {
				continue
			}

			path := msg.Path
			pathLabel.SetText(fmt.Sprintf("Path: %s", path))

			videoData := msg.GetServerVideoChunk().Data
			imgReader := bytes.NewReader(videoData)
			frame, err := jpeg.Decode(imgReader)
			if err != nil {
				log.Printf("Error decoding JPEG: %s\n", err)
				continue
			}

			bufferMutex.Lock()
			frameBuffer = append(frameBuffer, &frame)
			bufferMutex.Unlock()

			// Jitter calculation
			receivedTime := time.Now().UnixMilli()
			sendTimestamp := msg.GetServerVideoChunk().Timestamp
			currentDelay := float64(receivedTime - sendTimestamp)

			statsMutex.Lock()
			if packetCount > 0 {
				jitter := math.Abs(currentDelay - prevDelay)
				jittersum += jitter
			}
			prevDelay = currentDelay
			packetCount++
			totalSent = int(msg.GetServerVideoChunk().GetSequenceNumber())
			totalBytesReceived += int64(len(videoData))
			statsMutex.Unlock()
		}
	}()

	// Playback loop for frames in the buffer
	go func() {
		ticker := time.NewTicker(33 * time.Millisecond) // Approx 30 FPS
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				bufferMutex.Lock()
				if len(frameBuffer) > 0 {
					frame := *frameBuffer[0]
					frameBuffer = frameBuffer[1:]

					bufferMutex.Unlock()

					imgCanvas.Image = frame
					imgCanvas.Refresh()

					statsMutex.Lock()
					frameDisplayCount++
					elapsed := time.Since(lastFpsUpdate).Seconds()
					if elapsed >= 1.0 {
						fps = float64(frameDisplayCount) / elapsed
						frameDisplayCount = 0
						lastFpsUpdate = time.Now()
					}
					statsMutex.Unlock()
				} else {
					bufferMutex.Unlock()
				}
			}
		}
	}()

	myWindow.ShowAndRun()
}

// Helper function to send commands
func sendCommand(command, video, target string, conn net.Conn, config *cnf.AppConfigList, listenIp string, resetStats func()) {
	resetStats()

	message := protobuf.Header{
		Type:           protobuf.RequestType_RETRANSMIT,
		Length:         0,
		Timestamp:      time.Now().UnixMilli(),
		Sender:         "client",
		Path:           config.NodeName,
		Target:         []string{target},
		RequestedVideo: fmt.Sprintf("%s.Mjpeg", video),
		Content: &protobuf.Header_ClientCommand{
			ClientCommand: &protobuf.ClientCommand{
				Command:               protobuf.PlayerCommand(protobuf.PlayerCommand_value[command]),
				AdditionalInformation: "no additional information",
			},
		},
	}
	message.Length = int32(proto.Size(&message))
	data, err := proto.Marshal(&message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}
	_, err = conn.Write(data)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

// Show statistics in a new window
func showStatsWindow(app fyne.App, startTime time.Time, packetCount *int, jittersum *float64, statsMutex *sync.Mutex, totalSent *int, totalBytesReceived *int64, fps *float64) {
	statsWindow := app.NewWindow("Statistics")
	statsWindow.Resize(fyne.NewSize(400, 300))

	// Labels to display the stats
	elapsedTimeLabel := widget.NewLabel("Elapsed Time: 0.00 s")
	packetCountLabel := widget.NewLabel(fmt.Sprintf("Packets Received: %d", *packetCount))
	totalBytesLabel := widget.NewLabel(fmt.Sprintf("Total Bytes: %d", *totalBytesReceived))
	avgJitterLabel := widget.NewLabel(fmt.Sprintf("Avg Jitter: %.2f ms", *jittersum/float64(*packetCount)))
	fpsLabel := widget.NewLabel(fmt.Sprintf("FPS: %.2f", *fps))

	statsWindow.SetContent(container.NewVBox(elapsedTimeLabel, packetCountLabel, totalBytesLabel, avgJitterLabel, fpsLabel))

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			statsMutex.Lock()
			elapsed := time.Since(startTime).Seconds()
			elapsedTimeLabel.SetText(fmt.Sprintf("Elapsed Time: %.2f s", elapsed))
			packetCountLabel.SetText(fmt.Sprintf("Packets Received: %d", *packetCount))
			totalBytesLabel.SetText(fmt.Sprintf("Total Bytes: %d", *totalBytesReceived))
			avgJitter := 0.0
			if *packetCount > 0 {
				avgJitter = *jittersum / float64(*packetCount)
			}
			avgJitterLabel.SetText(fmt.Sprintf("Avg Jitter: %.2f ms", avgJitter))
			fpsLabel.SetText(fmt.Sprintf("FPS: %.2f", *fps))
			statsMutex.Unlock()
		}
	}()

	statsWindow.Show()
}
