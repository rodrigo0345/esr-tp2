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
	prevDelay   float64
	jitterSum   float64
	packetCount int
	totalSent   int
	startTime   time.Time
	mu          sync.Mutex
}

func (s *Statistics) UpdateStats(receivedTime int64, sendTimestamp int64, sequenceNumber int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentDelay := float64(receivedTime - sendTimestamp)
	if s.packetCount > 0 {
		jitter := math.Abs(currentDelay - s.prevDelay)
		s.jitterSum += jitter
	}
	s.prevDelay = currentDelay
	s.packetCount++
	s.totalSent = sequenceNumber
}

func (s *Statistics) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prevDelay = 0
	s.jitterSum = 0
	s.packetCount = 0
	s.totalSent = 0
	s.startTime = time.Now()
}

func (s *Statistics) CalculateMetrics() (float64, int, float64, float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	elapsedTime := time.Since(s.startTime).Seconds()
	avgJitter := 0.0
	if s.packetCount > 1 {
		avgJitter = s.jitterSum / float64(s.packetCount-1)
	}
	packetLoss := 0.0
	if s.totalSent > 0 {
		lossCount := s.totalSent - s.packetCount
		packetLoss = float64(lossCount) / float64(s.totalSent) * 100
	}
	return elapsedTime, s.packetCount, avgJitter, packetLoss
}

type UI struct {
	imgCanvas    *canvas.Image
	pathLabel    *widget.Label
	messageEntry *widget.Entry
	targetEntry  *widget.Entry
	sendButton   *widget.Button
	cancelButton *widget.Button
	statsButton  *widget.Button
	stats        *Statistics
}

func NewUI() *UI {
	return &UI{
		imgCanvas:    canvas.NewImageFromImage(nil),
		pathLabel:    widget.NewLabel(""),
		messageEntry: widget.NewEntry(),
		targetEntry:  widget.NewEntry(),
		sendButton:   widget.NewButton("Send", nil),
		cancelButton: widget.NewButton("Stop", nil),
		statsButton:  widget.NewButton("Statistics", nil),
		stats:        &Statistics{startTime: time.Now()},
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

	done := make(chan struct{})
	statsWindow.SetOnClosed(func() { close(done) })

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				elapsedTime, packetCount, avgJitter, packetLoss := ui.stats.CalculateMetrics()
				elapsedTimeLabel.SetText(fmt.Sprintf("Elapsed Time: %.2f seconds", elapsedTime))
				packetCountLabel.SetText(fmt.Sprintf("Packets Received: %d", packetCount))
				avgJitterLabel.SetText(fmt.Sprintf("Average Jitter: %.2f ms", avgJitter))
				packetLossRate.SetText(fmt.Sprintf("Packet Loss: %.2f%%", packetLoss))
				time.Sleep(500 * time.Millisecond)
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

	ui.messageEntry.SetPlaceHolder("Enter video")
	ui.targetEntry.SetPlaceHolder("Enter target")

	// Button actions
	ui.sendButton.OnTapped = func() {
		sendCommand("PLAY", ui.messageEntry.Text, ui.targetEntry.Text, bestPopAddr, cnf, ui.stats)
	}
	ui.cancelButton.OnTapped = func() {
		sendCommand("STOP", ui.messageEntry.Text, ui.targetEntry.Text, bestPopAddr, cnf, ui.stats)
	}
	ui.statsButton.OnTapped = func() {
		ui.ShowStatsWindow(myApp)
	}

	// Arrange buttons horizontally
	buttons := container.NewHBox(ui.sendButton, ui.cancelButton, ui.statsButton)

	// Place other controls in a vertical container
	controls := container.NewVBox(ui.messageEntry, ui.targetEntry, buttons, ui.pathLabel)

	// Combine the controls with the image canvas
	content := container.NewBorder(nil, controls, nil, nil, ui.imgCanvas)

	myWindow.SetContent(content)

	// Goroutine to process messages from uiChannel
	go func() {
		for msg := range uiChannel {
      if len(msg.Path) > 0 {
        ui.pathLabel.SetText(fmt.Sprintf("Path: %s", msg.Path))
      } else if len(msg.RequestedVideo) > 8 {
        ui.pathLabel.SetText(fmt.Sprintf("Path: %s", msg.Path[0:8]))
      }

			if videoChunk := msg.GetServerVideoChunk(); videoChunk != nil {
				ui.UpdateImage(videoChunk.Data)
				ui.stats.UpdateStats(time.Now().UnixMilli(), videoChunk.Timestamp, int(videoChunk.SequenceNumber))
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
