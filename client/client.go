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

	cnf "github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"

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
	defer listener.Close()

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
		totalSent          int
		totalBytesReceived int64
		statsMutex         sync.Mutex
	)

	resetStats := func() {
		statsMutex.Lock()
		defer statsMutex.Unlock()

		// Reset statistics values
		packetCount = 0
		jittersum = 0.0
		prevDelay = 0.0
		totalBytesReceived = 0
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
		showStatsWindow(myApp, &statsMutex, &packetCount, &jittersum, &totalSent, &totalBytesReceived)
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
			totalSent = int(msg.GetServerVideoChunk().GetSequenceNumber())
			totalBytesReceived += int64(len(videoData))
			statsMutex.Unlock()
		}
	}()

	myWindow.ShowAndRun()
}

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
func showStatsWindow(app fyne.App, statsMutex *sync.Mutex, packetCount *int, jittersum *float64, totalSent *int, totalBytesReceived *int64) {
	statsWindow := app.NewWindow("Statistics")
	statsWindow.Resize(fyne.NewSize(400, 300))

	elapsedTimeLabel := widget.NewLabel("")
	packetCountLabel := widget.NewLabel("")
	totalBytesLabel := widget.NewLabel("")
	avgJitterLabel := widget.NewLabel("")
	packetLossRate := widget.NewLabel("")

	statsContent := container.NewVBox(
		elapsedTimeLabel,
		packetCountLabel,
		totalBytesLabel,
		avgJitterLabel,
		packetLossRate,
	)
	statsWindow.SetContent(statsContent)

	// Update statistics periodically
	go func() {
		for range time.Tick(500 * time.Millisecond) {
			statsMutex.Lock()
			avgJitter := 0.0
			if *packetCount > 1 {
				avgJitter = *jittersum / float64(*packetCount-1)
			}
			packetLoss := 0.0
			if *totalSent > 0 {
				lostPackets := *totalSent - *packetCount
				packetLoss = (float64(lostPackets) / float64(*totalSent)) * 100
			}
			totalBytes := *totalBytesReceived
			statsMutex.Unlock()

			elapsedTimeLabel.SetText(fmt.Sprintf("Elapsed Time: %.2f seconds", time.Since(time.Now()).Seconds()))
			packetCountLabel.SetText(fmt.Sprintf("Packets Received: %d", *packetCount))
			totalBytesLabel.SetText(fmt.Sprintf("Total Bytes: %d", totalBytes))
			avgJitterLabel.SetText(fmt.Sprintf("Avg Jitter: %.2f ms", avgJitter))
			packetLossRate.SetText(fmt.Sprintf("Packet Loss: %.2f%%", packetLoss))
		}
	}()
	statsWindow.Show()
}
