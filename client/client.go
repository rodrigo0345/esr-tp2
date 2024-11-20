package client

import (
    "bytes"
    "fmt"
    "image"
    "image/jpeg"
    "log"
    "net"
    "time"

    "github.com/rodrigo0345/esr-tp2/config"
    "github.com/rodrigo0345/esr-tp2/config/protobuf" // Import your generated protobuf package
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

    listenIp := "127.0.0.1:2222"

    // Setup a UDP connection for sending
    conn, err := net.Dial("udp", pneIpString)
    if err != nil {
        log.Fatalf("Failed to dial: %v", err)
    }
    defer conn.Close()

    // Create a UDP listener for incoming messages
    laddr := net.UDPAddr{
        Port: 2222,
        IP:   net.ParseIP("127.0.0.1"),
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
    img := canvas.NewImageFromImage(nil)
    img.FillMode = canvas.ImageFillContain

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

    controls := container.NewVBox(messageEntry, targetEntry, sendButton)
    content := container.NewBorder(nil, controls, nil, nil, img)

    myWindow.SetContent(content)
    myWindow.Resize(fyne.NewSize(640, 480))

    // Channel to send images from the goroutine to the main thread
    imgChan := make(chan image.Image)

    // Run a goroutine to listen for incoming messages
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

            // Decode the JPEG data into an image.Image
            imgReader := bytes.NewReader(videoData)
            frame, err := jpeg.Decode(imgReader)
            if err != nil {
                log.Printf("Error decoding JPEG: %s\n", err)
                continue
            }

            // Send the image to the main thread via the channel
            imgChan <- frame
        }
    }()

    // In the main thread, receive images from the channel and update the image widget
        for frame := range imgChan {
            // Need to update the UI in the main thread
            fyne.CurrentApp().RunOnMain(func() {
                img.Image = frame
                img.Refresh()
            })
        }

    myWindow.ShowAndRun()
}

