package server

import (
	"fmt"
	"github.com/pion/rtp"
	"math/rand"
	"net"
	"time"
)

type Client struct {
	Source            *VideoStream
	RtpPort           int
	RtpAddr           string
	OtherThreadSignal chan bool
}

func (c *Client) NewClient(filename string, rtpPort int, rtpAddr string) {
	var err error
	c.Source, err = NewVideoStream(filename)
	if err != nil {
		panic(err)
	}

	c.RtpPort = rtpPort
	c.RtpAddr = rtpAddr
	c.OtherThreadSignal = make(chan bool)
}

// SendSource sends the video stream and updates the server-side display
func (c *Client) SendSource() {
	raddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", c.RtpAddr, c.RtpPort))
	if err != nil {
		panic(err)
	}

	socket, err := net.DialUDP("udp4", nil, raddr)
	if err != nil {
		panic(err)
	}
	defer socket.Close()

	// Increase the socket write buffer size (e.g., to 1MB)
	err = socket.SetWriteBuffer(1048576)

	fmt.Println("Sending source to: ", c.RtpAddr, c.RtpPort)

	// Randomly generated SSRC for this RTP stream
	ssrc := rand.Uint32()

	// Define max payload size (typically 1200 bytes for RTP over UDP with MTU considerations)
	const maxPayloadSize = 1400

	for {
		// Get the next frame from the source
		frame, err := c.Source.NextFrame()
		if err != nil {
			panic(err)
		}

		// Fragment the frame into chunks if needed
		frameLen := len(frame)
		numberOfPackets := (frameLen + maxPayloadSize - 1) / maxPayloadSize // Total packets needed

		for i := 0; i < numberOfPackets; i++ {
			// Calculate the payload size for this packet
			start := i * maxPayloadSize
			end := start + maxPayloadSize
			if end > frameLen {
				end = frameLen
			}
			payload := frame[start:end]

			// Create a new RTP packet
			rtpPacket := &rtp.Packet{
				Header: rtp.Header{
					Version:        2,                                       // RTP version
					PayloadType:    26,                                      // Assuming 26 is the appropriate payload type for JPEG
					SequenceNumber: uint16(c.Source.FrameNumber() + i),      // Sequence number increases with each packet
					Timestamp:      uint32(time.Now().UnixNano() / 1000000), // Timestamp in milliseconds
					SSRC:           ssrc,
				},
				Payload: payload, // Payload is the chunk of the frame
			}

			// Set marker bit only for the last packet of the frame
			if i == numberOfPackets-1 {
				rtpPacket.Marker = true
			}

			// Serialize the RTP packet to raw bytes
			rawPacket, err := rtpPacket.Marshal()
			if err != nil {
				panic(err)
			}

			// Send the packet over UDP
			_, err = socket.Write(rawPacket)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Sent RTP packet for frame #%d, part %d/%d, size: %d bytes\n", c.Source.FrameNumber(), i+1, numberOfPackets, len(rawPacket))
		}

		// Wait for the next frame (33ms for ~30 FPS)
		time.Sleep(time.Millisecond * 33)
	}
}
