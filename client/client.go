package client

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/pion/rtp"
	"github.com/rodrigo0345/esr-tp2/config"
)

// use GUI
func Client(config *config.AppConfigList) {
	a := app.New()
	w := a.NewWindow("RTP Client")

	// Start listening for UDP packets in a separate goroutine
	go Listen(config, w)

	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}


// Listen listens for incoming RTP packets and processes them
func Listen(config *config.AppConfigList, window fyne.Window) {
	addr := fmt.Sprintf("%s:%d", config.ServerUrl.Url, config.ServerUrl.Port)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}
	defer conn.Close()

	log.Printf("Listening on %s", addr)

	buffer := make([]byte, 1500) // Assuming typical RTP packet size
	for {
		// Read incoming UDP packet
		n, _, err := conn.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}

		// Process the RTP packet
		processRTPPacket(buffer[:n], window)
	}
}

// processRTPPacket decodes the RTP packet and displays the image
// Global buffer to accumulate the RTP packets for one image
var imageBuffer bytes.Buffer

// processRTPPacket decodes the RTP packet and displays the image
func processRTPPacket(packet []byte, window fyne.Window) {
	// Unmarshal the RTP packet using the pion/rtp package
	rtpPacket := &rtp.Packet{}
	if err := rtpPacket.Unmarshal(packet); err != nil {
		log.Printf("Failed to unmarshal RTP packet: %v", err)
		return
	}

	// Assuming the payload is a JPEG image (PayloadType 26 is JPEG in RTP)
	if rtpPacket.PayloadType != 26 {
		log.Printf("Unexpected payload type: %d", rtpPacket.PayloadType)
		return
	}

	// Accumulate the payloads in a buffer
	imageBuffer.Write(rtpPacket.Payload)

	// Check if this is the last packet of the current frame (marker bit set)
	if rtpPacket.Marker {
		// Decode and display the image when we receive the last packet
		img, err := decodeImage(imageBuffer.Bytes())
		if err != nil {
			log.Printf("Failed to decode image: %v", err)
			imageBuffer.Reset() // Clear the buffer in case of an error
			return
		}

		displayImage(img, window)
		imageBuffer.Reset() // Clear the buffer for the next frame
	}
}

func decodeImage(payload []byte) (image.Image, error) {
	// Decode the JPEG image from the RTP payload
	img, err := jpeg.Decode(bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("Failed to decode JPEG image: %v", err)
	}
	return img, nil
}

// displayImage displays the decoded image on the Fyne window
func displayImage(img image.Image, window fyne.Window) {
	// Convert image to PNG format to be compatible with Fyne canvas
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		log.Printf("Error encoding PNG: %v", err)
		return
	}

	// Create an image canvas object
	imageResource := fyne.NewStaticResource("Image", buf.Bytes())
	imgCanvas := canvas.NewImageFromResource(imageResource)
	imgCanvas.FillMode = canvas.ImageFillContain

	// Update the window content with the new image
	window.SetContent(imgCanvas)
}
