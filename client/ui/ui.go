package ui

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"log"
	"math"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/rodrigo0345/esr-tp2/config"
	cnf "github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
  )

  type Statistics struct {
    mu                 sync.Mutex
    startTime          time.Time
    dataPoints         []statDataPoint // Rolling window of stats
    totalBytesReceived int64
    totalSent          int
    pathLabel          string
  }

  type statDataPoint struct {
    timestamp      time.Time
    packetSize     int
    sequenceNumber int
    receivedTime   int64
    sendTimestamp  int64
  }

  // Maximum time window for metrics
  const rollingWindowSize = 5 * time.Second

  func (s *Statistics) Reset() {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.startTime = time.Now()
    s.dataPoints = nil
    s.totalBytesReceived = 0
    s.totalSent = 0
  }

  func (s *Statistics) AddDataPoint(receivedTime int64, sendTimestamp int64, sequenceNumber int, packetSize int, path string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now()
    s.dataPoints = append(s.dataPoints, statDataPoint{
      timestamp:      now,
      packetSize:     packetSize,
      sequenceNumber: sequenceNumber,
      receivedTime:   receivedTime,
      sendTimestamp:  sendTimestamp,
    })
    s.pathLabel = path
    s.totalBytesReceived += int64(packetSize)
    s.totalSent = sequenceNumber

    // Remove old data points beyond the rolling window
    cutoff := now.Add(-rollingWindowSize)
    for len(s.dataPoints) > 0 && s.dataPoints[0].timestamp.Before(cutoff) {
      s.dataPoints = s.dataPoints[1:]
    }
  }

  func (s *Statistics) CalculateMetrics() (float64, int, float64, float64, float64, float64) {
    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now()
    cutoff := now.Add(-rollingWindowSize)
    validDataPoints := 0
    jitterSum := 0.0
    prevDelay := 0.0
    totalBytes := 0
    firstSequence := -1
    lastSequence := -1

    for _, point := range s.dataPoints {
      if point.timestamp.Before(cutoff) {
        continue
      }
      validDataPoints++
      totalBytes += point.packetSize

      if firstSequence == -1 || point.sequenceNumber < firstSequence {
        firstSequence = point.sequenceNumber
      }
      if lastSequence == -1 || point.sequenceNumber > lastSequence {
        lastSequence = point.sequenceNumber
      }

      // Delay and jitter
      currentDelay := float64(point.receivedTime - point.sendTimestamp)
      if lastSequence >= 0 {
        jitter := math.Abs(currentDelay - prevDelay)
        jitterSum += jitter
      }
      prevDelay = currentDelay
    }

    // Calculate packet loss
    totalPackets := lastSequence - firstSequence + 1
    packetLoss := 0.0
    if totalPackets > 0 {
      packetLoss = float64(totalPackets-validDataPoints) / float64(totalPackets) * 100
    }

    // Calculate throughput in kilobits per second
    throughput := float64(totalBytes) * 8 / rollingWindowSize.Seconds() / 1024.0 // Convert to kbps

    // Calculate FPS
    fps := float64(validDataPoints) / rollingWindowSize.Seconds()

    // Average jitter
    avgJitter := 0.0
    if validDataPoints > 1 {
      avgJitter = jitterSum / float64(validDataPoints-1)
    }

    // Elapsed time
    elapsedTime := time.Since(s.startTime).Seconds()

    return elapsedTime, validDataPoints, avgJitter, packetLoss, fps, throughput
  }

  type UI struct {
    imgCanvas     *canvas.Image
    pathLabel     *widget.Label
    videoDropdown *widget.Select
    targetEntry   *widget.Entry
    sendButton    *widget.Button
    cancelButton  *widget.Button
    statsButton   *widget.Button
    stats         *Statistics
  }

  func NewUI() *UI {
    return &UI{
      imgCanvas:     canvas.NewImageFromImage(nil),
      pathLabel:     widget.NewLabel(""),
      videoDropdown: widget.NewSelect([]string{"film", "demo"}, nil), // Dropdown menu
      targetEntry:   widget.NewEntry(),
      sendButton:    widget.NewButton("Send", nil),
      cancelButton:  widget.NewButton("Stop", nil),
      statsButton:   widget.NewButton("Statistics", nil),
      stats:         &Statistics{startTime: time.Now()},
    }
  }

  func (ui *UI) UpdateImage(data []byte) {
    imgReader := bytes.NewReader(data)
    frame, err := jpeg.Decode(imgReader)
    if err != nil {
      log.Printf("Error decoding JPEG: %s\n", err)
      return
    }
    ui.imgCanvas.Image = frame
    ui.imgCanvas.Refresh()
  }

  func (ui *UI) ShowStatsWindow(app fyne.App) {
	statsWindow := app.NewWindow("Statistics")
	statsWindow.Resize(fyne.NewSize(400, 300))

	elapsedTimeLabel := widget.NewLabel("Elapsed Time: 0.00 s")
	packetCountLabel := widget.NewLabel("Packets Received: 0")
	avgJitterLabel := widget.NewLabel("Average Jitter: 0.00 ms")
	packetLossLabel := widget.NewLabel("Packet Loss: 0.00%")
	fpsLabel := widget.NewLabel("FPS: 0.00")
	throughputLabel := widget.NewLabel("Throughput: 0.00 kbps")
  pathLabel := widget.NewLabel("Path: ")

	statsContent := container.NewVBox(
		elapsedTimeLabel,
		packetCountLabel,
		avgJitterLabel,
		packetLossLabel,
		fpsLabel,
		throughputLabel,
    pathLabel,
	)
	statsWindow.SetContent(statsContent)

	done := make(chan struct{})
	statsWindow.SetOnClosed(func() { close(done) })

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				elapsedTime, packetCount, avgJitter, packetLoss, fps, throughput := ui.stats.CalculateMetrics()
				elapsedTimeLabel.SetText(fmt.Sprintf("Elapsed Time: %.2f s", elapsedTime))
				packetCountLabel.SetText(fmt.Sprintf("Packets Received: %d", packetCount))
				avgJitterLabel.SetText(fmt.Sprintf("Average Jitter: %.2f ms", avgJitter))
				packetLossLabel.SetText(fmt.Sprintf("Packet Loss: %.2f%%", packetLoss))
				fpsLabel.SetText(fmt.Sprintf("FPS: %.2f", fps))
				throughputLabel.SetText(fmt.Sprintf("Throughput: %.2f kbps", throughput))
        pathLabel.SetText(fmt.Sprintf("Path: %s", ui.stats.pathLabel))
				time.Sleep(1200 * time.Millisecond)
			}
		}
	}()

	statsWindow.Show()
}

func StartUI(bestPopAddr string, cnf *cnf.AppConfigList, uiChannel <-chan *protobuf.Header) {
	ui := NewUI()
	myApp := app.New()
	myWindow := myApp.NewWindow("Streaming Client")

	// Set window width constraints
	const maxWidth = 600
	myWindow.Resize(fyne.NewSize(maxWidth, 600))
	myWindow.SetFixedSize(true) // Prevent resizing beyond the set width and height

	ui.videoDropdown.PlaceHolder = "Select video" // Placeholder text for dropdown
	ui.targetEntry.SetPlaceHolder("Enter target")

	// Button actions
	ui.sendButton.OnTapped = func() {
		if ui.videoDropdown.Selected == "" {
			log.Println("Please select a video before sending.")
			return
		}
		sendCommand("PLAY", ui.videoDropdown.Selected, ui.targetEntry.Text, bestPopAddr, cnf, ui.stats)
	}
	ui.cancelButton.OnTapped = func() {
		sendCommand("STOP", ui.videoDropdown.Selected, ui.targetEntry.Text, bestPopAddr, cnf, ui.stats)
	}
	ui.statsButton.OnTapped = func() {
		ui.ShowStatsWindow(myApp)
	}

	// Arrange buttons horizontally
	buttons := container.NewHBox(ui.sendButton, ui.cancelButton, ui.statsButton)

	// Place other controls in a vertical container
	controls := container.NewVBox(ui.videoDropdown, ui.targetEntry, buttons, ui.pathLabel)

	// Combine the controls with the image canvas
	content := container.NewBorder(nil, controls, nil, nil, ui.imgCanvas)

	myWindow.SetContent(content)

	// Goroutine to process messages from uiChannel
	go func() {
		for msg := range uiChannel {

			if videoChunk := msg.GetServerVideoChunk(); videoChunk != nil {
				ui.UpdateImage(videoChunk.Data)
				ui.stats.AddDataPoint(
					time.Now().UnixMilli(),
					videoChunk.Timestamp,
					int(videoChunk.SequenceNumber),
					len(videoChunk.Data),
          msg.Path,
				)
			}
		}
	}()

	myWindow.ShowAndRun()
}

func sendCommand(command, video, target string, bestPopAddr string, cnf *cnf.AppConfigList, stats *Statistics) {
	stats.Reset()

	message := protobuf.Header{
		Type:           protobuf.RequestType_RETRANSMIT,
		Length:         0,
		Timestamp:      time.Now().UnixMilli(),
		Sender:         "client",
		Path:           cnf.NodeName,
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
	config.SendMessageUDP(bestPopAddr, data)
}
