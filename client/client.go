package client

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"log"
	"math"
	"net"
	"sync"
	"time"

	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
  cnf "github.com/rodrigo0345/esr-tp2/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
		prevDelay   float64
		jittersum   float64
		packetCount int
		startTime   = time.Now()
		statsMutex  sync.Mutex
		totalSent   int // Total number of packets received (sequence number)
	)

	resetStats := func() {
		statsMutex.Lock()
		defer statsMutex.Unlock()

		// Reset statistics values
		packetCount = 0
		jittersum = 0.0
		prevDelay = 0.0
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
		showStatsWindow(myApp, startTime, &packetCount, &jittersum, &statsMutex, &totalSent)
	})

	controls := container.NewVBox(messageEntry, targetEntry, sendButton, cancelButton, statsButton, pathLabel)
	content := container.NewBorder(nil, controls, nil, nil, imgCanvas)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 600)) // Adjusted size of the main window

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

			imgCanvas.Image = frame
			imgCanvas.Refresh()

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
			totalSent = int(msg.GetServerVideoChunk().GetSequenceNumber()) // Update totalSent with the sequence number
			statsMutex.Unlock()
		}
	}()

	myWindow.ShowAndRun()
}

// Helper function to send commands
// Helper function to send commands
func sendCommand(command, video, target string, conn net.Conn, config *cnf.AppConfigList, listenIp string, resetStats func()) {
	// Reset statistics before sending the command
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
func showStatsWindow(app fyne.App, startTime time.Time, packetCount *int, jittersum *float64, statsMutex *sync.Mutex, totalSent *int) {
	statsWindow := app.NewWindow("Statistics")
	statsWindow.Resize(fyne.NewSize(400, 300))

	// Labels para exibição das estatísticas
	elapsedTimeLabel := widget.NewLabel("")
	packetCountLabel := widget.NewLabel("")
	avgJitterLabel := widget.NewLabel("")
	packetLossRate := widget.NewLabel("")

	statsContent := container.NewVBox(
		elapsedTimeLabel,
		packetCountLabel,
		avgJitterLabel,
		packetLossRate,
	)
	statsWindow.SetContent(statsContent)

	// Channel to signal window closure
	done := make(chan struct{})

	statsWindow.SetOnClosed(func() {
		close(done)
	})

	go func() {
		for {
			select {
			case <-done:
				return // Close Goroutine when the window is close
			default:
				statsMutex.Lock()
				// Elapsed time
				elapsedTime := time.Since(startTime).Seconds()

				// Calculate Avg Jitter
				avgJitter := 0.0
				if *packetCount > 1 {
					avgJitter = *jittersum / float64(*packetCount-1)
				}

				// Calculate Packet Loss
				expectedPackets := *totalSent
				lossCount := expectedPackets - *packetCount

				packetLossPercentage := 0.0
				if expectedPackets > 0 {
					packetLossPercentage = float64(lossCount) / float64(expectedPackets) * 100
				}

				statsMutex.Unlock()

				// Update the label values
				elapsedTimeLabel.SetText(fmt.Sprintf("Elapsed Time: %.2f seconds", elapsedTime))
				packetCountLabel.SetText(fmt.Sprintf("Packets Received: %d", *packetCount))
				avgJitterLabel.SetText(fmt.Sprintf("Average Jitter: %.2f ms", avgJitter))
				packetLossRate.SetText(fmt.Sprintf("Packet Loss: %.2f%%", packetLossPercentage))

				time.Sleep(500 * time.Millisecond) // Update every 500ms
			}
		}
	}()

	statsWindow.Show()
}
