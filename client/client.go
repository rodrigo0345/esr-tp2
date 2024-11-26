package client

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"log"
	"net"
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

type BestPop struct {
  Ip *protobuf.Interface
  Avg time.Duration
}

func Client(cnf *config.AppConfigList) {
	// Address to send messages
  logger := config.NewLogger(2, cnf.NodeName)

  pops, err := GetPops(cnf, logger)

  // TODO: ping all the pops and select the one with the lowest latency
  bestPop := GetBestPop(pops, logger)

	pneIpString := fmt.Sprintf("%s:%d", bestPop.Ip.Ip, bestPop.Ip.Port)
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

	// Create entry widgets and a button for sending messages
	messageEntry := widget.NewEntry()
	messageEntry.SetPlaceHolder("Enter message")
	targetEntry := widget.NewEntry()
	targetEntry.SetPlaceHolder("Enter target")
	sendButton := widget.NewButton("Send", func() {
		text := messageEntry.Text
		target := targetEntry.Text

		message := protobuf.Header{
			Type:           protobuf.RequestType_RETRANSMIT,
			Length:         0,
			Timestamp:      time.Now().UnixMilli(),
			ClientIp:       listenIp,
			Sender:         "client",
      Path:           cnf.NodeName,
			Target:         target,
			RequestedVideo: "lol.Mjpeg",
			Content: &protobuf.Header_ClientCommand{
				ClientCommand: &protobuf.ClientCommand{
					Command:               protobuf.PlayerCommand_PLAY,
					AdditionalInformation: text,
				},
			},
		}

		fmt.Printf("Sending message: %s to %s\n", text, target)
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
	})

	controls := container.NewVBox(messageEntry, targetEntry, sendButton, pathLabel)
	content := container.NewBorder(nil, controls, nil, nil, imgCanvas)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(640, 480))

	// Run a goroutine to listen for incoming messages
	fmt.Println("Listening for messages")
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

			path := msg.Path
			pathLabel.SetText(fmt.Sprintf("Path: %s", path))

			videoData := msg.GetServerVideoChunk().Data

			// Decode the JPEG data into an image.Image
			imgReader := bytes.NewReader(videoData)
			frame, err := jpeg.Decode(imgReader)
			if err != nil {
				log.Printf("Error decoding JPEG: %s\n", err)
				continue
			}

			// Update the image canvas with the new frame
			frameCopy := frame // Create a copy to avoid concurrency issues
			imgCanvas.Image = frameCopy
			imgCanvas.Refresh()
		}
	}()

	myWindow.ShowAndRun()
}
