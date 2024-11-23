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

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Client(config *config.AppConfigList) {
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

	// Buffer and statistics variables
	var (
		frameBuffer    = make([]image.Image, 0, 50) // Circular buffer for frames
		bufferMutex    sync.Mutex
		packetCount    int
		prevDelay      float64
		jitterSum      float64
		startTime      = time.Now()
		statsMutex     sync.Mutex
		totalSent      int
		packetLossRate float64
	)

	resetStats := func() {
		statsMutex.Lock()
		defer statsMutex.Unlock()

		packetCount = 0
		prevDelay = 0
		jitterSum = 0
		packetLossRate = 0
	}

	// Entry widgets and buttons
	messageEntry := widget.NewEntry()
	messageEntry.SetPlaceHolder("Enter video")
	targetEntry := widget.NewEntry()
	targetEntry.SetPlaceHolder("Enter target")
	sendButton := widget.NewButton("Send", func() {
		sendCommand("PLAY", messageEntry.Text, targetEntry.Text, conn, config, listenIp, resetStats)
	})
	cancelButton := widget.NewButton("Stop", func() {
		sendCommand("STOP", messageEntry.Text, targetEntry.Text, conn, config, listenIp, resetStats)
	})
	statsButton := widget.NewButton("Statistics", func() {
		showStatsWindow(myApp, startTime, &packetCount, &jitterSum, &statsMutex, &totalSent, &packetLossRate)
	})

	controls := container.NewVBox(messageEntry, targetEntry, sendButton, cancelButton, statsButton, pathLabel)
	content := container.NewBorder(nil, controls, nil, nil, imgCanvas)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 600))

	// Goroutine for listening to incoming messages
	go func() {
		for {
			buffer := make([]byte, 65535)
			n, _, err := listener.ReadFromUDP(buffer)

			if err != nil {
				log.Printf("Error receiving message: %s\n", err)
				continue
			}

			var msg protobuf.Header
			err = proto.Unmarshal(buffer[:n], &msg)
			if err != nil {
				log.Printf("Error unmarshaling message: %s\n", err)
				continue
			}

			if msg.GetServerVideoChunk() == nil {
				continue
			}

			videoData := msg.GetServerVideoChunk().Data
			imgReader := bytes.NewReader(videoData)
			frame, err := jpeg.Decode(imgReader)
			if err != nil {
				log.Printf("Error decoding JPEG: %s\n", err)
				continue
			}

			// Add frame to buffer
			bufferMutex.Lock()
			if len(frameBuffer) < cap(frameBuffer) {
				frameBuffer = append(frameBuffer, frame)
			} else {
				frameBuffer = frameBuffer[1:] // Remove oldest frame
				frameBuffer = append(frameBuffer, frame)
			}
			bufferMutex.Unlock()

			// Update path label
			path := msg.Path
			pathLabel.SetText(fmt.Sprintf("Path: %s", path))

			// Update statistics
			receivedTime := time.Now().UnixMilli()
			sendTimestamp := msg.GetServerVideoChunk().Timestamp
			currentDelay := float64(receivedTime - sendTimestamp)

			statsMutex.Lock()
			if packetCount > 0 {
				jitter := math.Abs(currentDelay - prevDelay)
				jitterSum += jitter
			}
			prevDelay = currentDelay
			packetCount++
			totalSent = int(msg.GetServerVideoChunk().GetSequenceNumber()) // Update totalSent with the sequence number
			if totalSent > 0 {
				lossCount := totalSent - packetCount
				packetLossRate = (float64(lossCount) / float64(totalSent)) * 100
			}
			statsMutex.Unlock()
		}
	}()

	// Playback loop for frames in the buffer
	go func() {
    for {
        time.Sleep(33 * time.Millisecond) // Approx 30 FPS

        bufferMutex.Lock()
        if len(frameBuffer) > 5 { // Start playback only if there are at least 5 frames
            frame := frameBuffer[0]
            if len(frameBuffer) > 1 {
                frameBuffer = frameBuffer[1:]
            } else {
                frameBuffer = frameBuffer[:0]
            }
            bufferMutex.Unlock()

            imgCanvas.Image = frame
            imgCanvas.Refresh()
        } else {
            bufferMutex.Unlock()
        }
    }
  }()

	myWindow.ShowAndRun()
}

// Statistics Window
func showStatsWindow(app fyne.App, startTime time.Time, packetCount *int, jitterSum *float64, statsMutex *sync.Mutex, totalSent *int, packetLossRate *float64) {
	statsWindow := app.NewWindow("Statistics")
	statsWindow.Resize(fyne.NewSize(400, 300))

	elapsedTimeLabel := widget.NewLabel("")
	packetCountLabel := widget.NewLabel("")
	avgJitterLabel := widget.NewLabel("")
	packetLossLabel := widget.NewLabel("")

	statsContent := container.NewVBox(
		elapsedTimeLabel,
		packetCountLabel,
		avgJitterLabel,
		packetLossLabel,
	)
	statsWindow.SetContent(statsContent)

	done := make(chan struct{})

	statsWindow.SetOnClosed(func() {
		close(done)
	})

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				statsMutex.Lock()
				elapsedTime := time.Since(startTime).Seconds()
				avgJitter := 0.0
				if *packetCount > 1 {
					avgJitter = *jitterSum / float64(*packetCount-1)
				}
				statsMutex.Unlock()

				// Update labels
				elapsedTimeLabel.SetText(fmt.Sprintf("Elapsed Time: %.2f seconds", elapsedTime))
				packetCountLabel.SetText(fmt.Sprintf("Packets Received: %d", *packetCount))
				avgJitterLabel.SetText(fmt.Sprintf("Average Jitter: %.2f ms", avgJitter))
				packetLossLabel.SetText(fmt.Sprintf("Packet Loss Rate: %.2f%%", *packetLossRate))

				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	statsWindow.Show()
}

// SendCommand Function
func sendCommand(command, video, target string, conn net.Conn, config *config.AppConfigList, listenIp string, resetStats func()) {
	// Reset statistics before sending the command
	resetStats()

	// Construct the protobuf message
	message := protobuf.Header{
		Type:           protobuf.RequestType_RETRANSMIT,
		Length:         0,
		Timestamp:      time.Now().UnixMilli(),
		ClientPort:     string(2222),
		Sender:         "client",
		Path:           config.NodeName,
		Target:         []string{target}, // Updated to correctly handle the target as a slice
		RequestedVideo: fmt.Sprintf("%s.Mjpeg", video),
		Content: &protobuf.Header_ClientCommand{
			ClientCommand: &protobuf.ClientCommand{
				Command:               protobuf.PlayerCommand(protobuf.PlayerCommand_value[command]),
				AdditionalInformation: "no additional information",
			},
		},
	}

	// Calculate the message length
	message.Length = int32(proto.Size(&message))

	// Marshal the message to send it over the network
	data, err := proto.Marshal(&message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	// Send the message via the UDP connection
	_, err = conn.Write(data)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

