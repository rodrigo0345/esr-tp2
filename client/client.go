package client

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
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
	laddr, err := net.ResolveUDPAddr("udp", listenIp)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	listener, err := net.ListenUDP("udp", laddr)
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

	// Create a label to display the msg.Path with fixed width
	pathLabel := widget.NewLabel("")
	pathLabel.Wrapping = 0

	// Set fixed size for the path label
	pathLabel.Resize(fyne.NewSize(300, 20))

	// Create a container with fixed width for the path label
	fixedPathContainer := container.NewHBox(
		widget.NewLabel("Path:"),
		pathLabel,
	)

	// Buffer and statistics variables
	var (
		frameBuffer        = make([]image.Image, 0, 50) // Circular buffer for frames
		bufferMutex        sync.Mutex
		packetCount        int
		prevDelay          float64
		jitterSum          float64
		startTime          = time.Now()
		statsMutex         sync.Mutex
		totalSent          int
		packetLossRate     float64
		totalBytesReceived int64
		errorCount         int
		frameDisplayCount  int
		fps                float64
		lastFpsUpdate      time.Time = time.Now()
	)

	resetStats := func() {
		statsMutex.Lock()
		defer statsMutex.Unlock()

		packetCount = 0
		prevDelay = 0
		jitterSum = 0
		packetLossRate = 0
		totalBytesReceived = 0
		errorCount = 0
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
		showStatsWindow(
			myApp,
			startTime,
			&packetCount,
			&jitterSum,
			&statsMutex,
			&totalSent,
			&packetLossRate,
			&totalBytesReceived,
			&errorCount,
			&fps,
		)
	})

	controls := container.NewVBox(
		messageEntry,
		targetEntry,
		sendButton,
		cancelButton,
		statsButton,
		fixedPathContainer, // Use the fixed path container
	)
	content := container.NewBorder(nil, controls, nil, nil, imgCanvas)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 600))

	// Goroutine for listening to incoming messages
	go func() {
		var partialData []byte
		for {
			buffer := make([]byte, 65535)
			n, _, err := listener.ReadFromUDP(buffer)
			if err != nil {
				log.Printf("Error receiving message: %s\n", err)
				continue
			}

			partialData = append(partialData, buffer[:n]...) // Append received data

			var msg protobuf.Header
			err = proto.Unmarshal(partialData, &msg)
			if err != nil {
				// Message is incomplete, wait for more data
				continue
			}

			// Successfully unmarshaled; process the message
			partialData = nil // Reset buffer for the next message
			handleIncomingMessage(&msg, &frameBuffer, &bufferMutex, &statsMutex, &packetCount, &prevDelay, &jitterSum, &totalBytesReceived, &totalSent, &packetLossRate, &errorCount, pathLabel)
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
					frame := frameBuffer[0]
					if len(frameBuffer) > 1 {
						frameBuffer = frameBuffer[1:]
					} else {
						frameBuffer = frameBuffer[:0]
					}
					bufferMutex.Unlock()

					imgCanvas.Image = frame
					imgCanvas.Refresh()

					// Update FPS
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

// Statistics Window
func showStatsWindow(app fyne.App, startTime time.Time, packetCount *int, jitterSum *float64, statsMutex *sync.Mutex, totalSent *int, packetLossRate *float64, totalBytesReceived *int64, errorCount *int, fps *float64) {
	statsWindow := app.NewWindow("Statistics")
	statsWindow.Resize(fyne.NewSize(400, 400))

	elapsedTimeLabel := widget.NewLabel("")
	packetCountLabel := widget.NewLabel("")
	avgJitterLabel := widget.NewLabel("")
	packetLossLabel := widget.NewLabel("")
	totalBytesLabel := widget.NewLabel("")
	errorRateLabel := widget.NewLabel("")
	fpsLabel := widget.NewLabel("")

	// Creating progress bars for visual representation
	jitterProgress := widget.NewProgressBar()
	jitterProgress.Min = 0
	jitterProgress.Max = 100 // Adjust as needed based on expected jitter
	packetLossProgress := widget.NewProgressBar()
	packetLossProgress.Min = 0
	packetLossProgress.Max = 100
	errorRateProgress := widget.NewProgressBar()
	errorRateProgress.Min = 0
	errorRateProgress.Max = 100

	// Labels for progress bars
	jitterLabel := widget.NewLabel("Average Jitter (ms):")
	packetLossLabelTitle := widget.NewLabel("Packet Loss Rate (%):")
	errorRateLabelTitle := widget.NewLabel("Error Rate (%):")

	statsContent := container.NewVBox(
		elapsedTimeLabel,
		packetCountLabel,
		avgJitterLabel,
		jitterLabel,
		jitterProgress,
		packetLossLabelTitle,
		packetLossProgress,
		errorRateLabelTitle,
		errorRateProgress,
		packetLossLabel,
		totalBytesLabel,
		errorRateLabel,
		fpsLabel,
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
				currentFps := *fps
				totalBytes := *totalBytesReceived
				errors := *errorCount
				packetLoss := *packetLossRate
				statsMutex.Unlock()

				// Update labels
				elapsedTimeLabel.SetText(fmt.Sprintf("Elapsed Time: %.2f seconds", elapsedTime))
				packetCountLabel.SetText(fmt.Sprintf("Packets Received: %d", *packetCount))
				avgJitterLabel.SetText(fmt.Sprintf("Average Jitter: %.2f ms", avgJitter))
				packetLossLabel.SetText(fmt.Sprintf("Packet Loss Rate: %.2f%%", packetLoss))
				totalBytesLabel.SetText(fmt.Sprintf("Total Bytes Received: %d bytes", totalBytes))
				errorRateLabel.SetText(fmt.Sprintf("Error Rate: %d errors", errors))
				fpsLabel.SetText(fmt.Sprintf("Frame Rate (FPS): %.2f", currentFps))

				// Update progress bars
				jitterProgress.SetValue(avgJitter)
				if avgJitter > 100 {
					jitterProgress.SetValue(100)
				}

				packetLossProgress.SetValue(packetLoss)
				if packetLoss > 100 {
					packetLossProgress.SetValue(100)
				}

				// Assuming errorCount represents total errors; calculate error rate based on totalSent
				var errorRate float64
				if *totalSent > 0 {
					errorRate = (float64(*errorCount) / float64(*totalSent)) * 100
				} else {
					errorRate = 0
				}
				errorRateProgress.SetValue(errorRate)
				if errorRate > 100 {
					errorRateProgress.SetValue(100)
				}

				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	statsWindow.Show()
}

func handleIncomingMessage(msg *protobuf.Header, frameBuffer *[]image.Image, bufferMutex, statsMutex *sync.Mutex, packetCount *int, prevDelay *float64, jitterSum *float64, totalBytesReceived *int64, totalSent *int, packetLossRate *float64, errorCount *int, pathLabel *widget.Label) {
	if msg.GetServerVideoChunk() == nil {
		return
	}

	videoData := msg.GetServerVideoChunk().Data
	imgReader := bytes.NewReader(videoData)
	frame, err := jpeg.Decode(imgReader)
	if err != nil {
		log.Printf("Error decoding JPEG: %s\n", err)
		statsMutex.Lock()
		(*errorCount)++
		statsMutex.Unlock()
		return
	}

	bufferMutex.Lock()
	*frameBuffer = append(*frameBuffer, frame)
	bufferMutex.Unlock()

	pathLabel.SetText(msg.Path)
	statsMutex.Lock()
	*packetCount++
	*totalBytesReceived += int64(len(videoData))
	*totalSent = int(msg.GetServerVideoChunk().GetSequenceNumber())
	statsMutex.Unlock()
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
		ClientPort:     "2222", // Correctly set as string
		Sender:         "client",
		Path:           config.NodeName,
		Target:         []string{target}, // Correctly handle the target as a slice
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

// Helper function to read a frame from the ffmpeg output
func readFrame(reader *bufio.Reader) ([]byte, error) {
	// JPEG SOI and EOI markers
	const (
		SOIMarker = 0xFFD8
		EOIMarker = 0xFFD9
	)

	var frameData []byte

	// Read until SOI marker is found
	for {
		b1, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if b1 == 0xFF {
			b2, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}
			if b2 == 0xD8 {
				// Found SOI marker
				frameData = append(frameData, 0xFF, 0xD8)
				break
			}
		}
	}

	// Read until EOI marker is found
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		frameData = append(frameData, b)
		if b == 0xD9 && len(frameData) >= 2 && frameData[len(frameData)-2] == 0xFF {
			// Found EOI marker
			break
		}
	}

	return frameData, nil
}
