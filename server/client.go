package server

import (
	"fmt"
	"net"
	"time"
)

type VideoEvent int

const (
  Start VideoEvent = iota
  Stop
  Pause
  Resume
)

type Client struct {
	Source            *VideoStream
	RtpPort           int
	RtpAddr           string
	OtherThreadSignal chan bool
}

func (c* Client) NewClient(filename string, rtpPort int, rtpAddr string) {
  var err error
	c.Source, err = NewVideoStream(filename)
  if err != nil {
    panic(err)
  }

	c.RtpPort = rtpPort
	c.RtpAddr = rtpAddr
	c.OtherThreadSignal = make(chan bool)
}

func (c *Client) SendSource() {
	// send source
  socket, err := net.Dial("udp4", fmt.Sprintf("%s:%d", c.RtpAddr, c.RtpPort))
  if err != nil {
    panic(err)
  }
  defer socket.Close()

  fmt.Println("Sending source to: ", c.RtpAddr, c.RtpPort)
  start := true

  for {
    // get the source next frame
    frame, err := c.Source.NextFrame()
    fmt.Println("Frame: ", c.Source.FrameNumber())

    if err != nil {
      panic(err)
    }

    // send frame
    rtpData := MakeRtpPacket(frame, uint16(c.Source.FrameNumber()), start)
    start = false

    // join header and payload
    packet := rtpData.GetPacket()
    _, err = socket.Write(packet)
    if err != nil {
      panic(err)
    }

    time.Sleep(time.Millisecond * 33)
  }

}
