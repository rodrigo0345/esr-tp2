package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"github.com/rodrigo0345/esr-tp2/presence"
	"github.com/rodrigo0345/esr-tp2/server/boostrapper"
)

type ServerSystem struct {
	PresenceSystem *presence.PresenceSystem
	Logger         *config.Logger
	VideoStreams   VideoStreams
	Bootstrapper   *boostrapper.Bootstrapper
	StreamNotify   chan StreamNotification
}

type StreamNotification struct {
	Action string // "start" or "stop"
	Video  *Stream
}

func NewServerSystem(cnf *config.AppConfigList) *ServerSystem {
	return &ServerSystem{
		PresenceSystem: presence.NewPresenceSystem(cnf),
		Logger:         config.NewLogger(2, cnf.NodeName),
		Bootstrapper:   boostrapper.NewBootstrapper("./server/boostrapper/nb.json", config.NewLogger(2, cnf.NodeName)),
		VideoStreams:   *NewVideoStreams(),
		StreamNotify:   make(chan StreamNotification, 10),
	}
}

func (ss *ServerSystem) HeartBeatNeighbors(seconds int) {
	ss.PresenceSystem.HeartBeatNeighbors(seconds)
}

func (ss *ServerSystem) ListenForClients() {
	// customize for the server
	tls := config.GenerateTLS()
	port := fmt.Sprintf(":%d", ss.PresenceSystem.Config.NodeIP.Port)

	listener, err := quic.ListenAddr(port, tls, nil)
	if err != nil {
		ss.Logger.Error(err.Error())
		return
	}

	ss.Logger.Info(fmt.Sprintf("QUIC server is listening on port %d\n", ss.PresenceSystem.Config.NodeIP.Port))

	for {
		// Accept incoming QUIC connections
		session, err := listener.Accept(context.Background())
		if err != nil {
			ss.Logger.Error(err.Error())
		}

		// For each session, handle incoming streams
		go func() {
			stream, err := session.AcceptStream(context.Background())
			if err != nil {
				ss.Logger.Error(err.Error())
				return
			}

			go func(stream quic.Stream) {
				msg, err := config.ReceiveMessage(stream)
				if err != nil {
					ss.Logger.Error(err.Error())
					return
				}

        ss.Logger.Info(fmt.Sprintf("Received message \n"))

				// unmarshal protobuf
				header, err := config.UnmarshalHeader(msg)

				if err != nil {
					ss.Logger.Error(err.Error())
					return
				}

				switch header.Type {
				case protobuf.RequestType_ROUTINGTABLE:

					presence.HandleRouting(ss.PresenceSystem, session, stream, header)
					break

				case protobuf.RequestType_BOOTSTRAPER:

					ss.Bootstrapper.Bootstrap(session, stream, header)
					break

				case protobuf.RequestType_RETRANSMIT:

					// check if the message is for this server
					/* if header.GetTarget() != ss.PresenceSystem.Config.NodeName {
						ss.Logger.Error(fmt.Sprintf("Received message for %s, but this is %s\n", header.GetTarget(), ss.PresenceSystem.Config.NodeName))
						ss.Logger.Info(fmt.Sprintf("Trying to retransmit to %s\n", header.GetTarget()))

						presence.HandleRetransmitFromClient(ss.PresenceSystem, header)
						return
					}*/

					if header.GetClientCommand() == nil {
						return
					}

					// this specifies where the message needs to go
					sender := make([]string, 1)
					sender[0] = header.Sender
					header.Target = sender
					header.Sender = ss.PresenceSystem.Config.NodeName
					stream.Close()

					switch header.GetClientCommand().Command {
					case protobuf.PlayerCommand_PLAY:
						ss.Logger.Info("Starting video stream")
						ss.StartVideoStream(header)
						break
					case protobuf.PlayerCommand_STOP:
						ss.Logger.Info("Stopping video stream")
						ss.StopVideoStream(header)
						break
					default:
						ss.Logger.Error("Unknown command")
						ss.StopVideoStream(header)
						break
					}
				default:
					break
				}

			}(stream)
		}()
	}
}

func (ss *ServerSystem) StartVideoStream(header *protobuf.Header) {
	videoChoice := header.RequestedVideo
	clientName := header.GetTarget()[0]
	ss.Logger.Info(fmt.Sprintf("Starting video stream for %s, requesting %s\n", clientName, videoChoice))

	client := &Client{PresenceNodeName: clientName, ClientIP: header.GetSender()}
	videoStream := ss.VideoStreams.AddStream(videoChoice, client)

	// Notify the background routine
	ss.StreamNotify <- StreamNotification{Action: "start", Video: videoStream}
}

func (ss *ServerSystem) StopVideoStream(header *protobuf.Header) {
	videoChoice := header.RequestedVideo
	clientName := header.GetTarget()[0]

	client := &Client{PresenceNodeName: clientName, ClientIP: header.GetSender()}
	videoStream := ss.VideoStreams.RemoveStream(videoChoice, *client)

	// Notify the background routine
	if videoStream != nil {
		ss.StreamNotify <- StreamNotification{Action: "stop", Video: videoStream}
	}
}

func (ss *ServerSystem) BackgroundStreaming() {
	activeStreams := make(map[string]chan struct{}) // To manage active streams and stop signals

	for {
		select {
		case notification := <-ss.StreamNotify:
			switch notification.Action {
			case "start":
				if _, exists := activeStreams[notification.Video.Video]; !exists {
					stopChan := make(chan struct{})
					activeStreams[notification.Video.Video] = stopChan
					go ss.streamVideo(notification.Video, stopChan)
				}
			case "stop":
				if stopChan, exists := activeStreams[notification.Video.Video]; exists {
					close(stopChan) // Signal the streaming routine to stop
					delete(activeStreams, notification.Video.Video)
				}
			}
		}
	}
}

func (ss *ServerSystem) streamVideo(video *Stream, stopChan chan struct{}) {
	pwd := os.Getenv("PWD")
	videoFilePath := fmt.Sprintf("%s/videos/%s", pwd, video.Video)

	ss.Logger.Info(fmt.Sprintf("Starting stream for video: %s, file path: %s", video.Video, videoFilePath))

	for {
		cmd := exec.Command("ffmpeg", "-stream_loop", "-1", "-i", videoFilePath, "-f", "image2pipe", "-vcodec", "mjpeg", "-")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			ss.Logger.Error(fmt.Sprintf("Failed to start ffmpeg process for video %s: %v", video.Video, err))
			return
		}

		if err := cmd.Start(); err != nil {
			ss.Logger.Error(fmt.Sprintf("Failed to execute ffmpeg for video %s: %v", video.Video, err))
			return
		}

		reader := bufio.NewReader(stdout)
		defer cmd.Wait()

		sqNumber := 0

		for {
			select {
			case <-stopChan:
				cmd.Process.Kill() // Terminate ffmpeg process
				ss.Logger.Info(fmt.Sprintf("Stream stopped for video %s", video.Video))
				return
			default:
				frameData, err := readFrame(reader)
				if err != nil {
					if err == io.EOF {
						ss.Logger.Info(fmt.Sprintf("End of video stream for video %s, restarting loop...", video.Video))
						break // Exit the loop and restart ffmpeg
					}
					ss.Logger.Error(fmt.Sprintf("Error reading frame for video %s: %v", video.Video, err))
					continue
				}

				sqNumber += 1
				var targets []string
				for _, client := range video.Clients {
					targets = append(targets, client.PresenceNodeName)
				}

				ss.Logger.Info(fmt.Sprintf("Sending video chunk to %s\n", targets))

				header := &protobuf.Header{
					Sender:         ss.PresenceSystem.Config.NodeName,
					Target:         targets,
					RequestedVideo: video.Video,
					Path:           ss.PresenceSystem.Config.NodeName,
					Type:           protobuf.RequestType_RETRANSMIT,
					Length:         int32(len(frameData)),
					Timestamp:      time.Now().UnixMilli(),
					Content: &protobuf.Header_ServerVideoChunk{
						ServerVideoChunk: &protobuf.ServerVideoChunk{
							Data:           frameData,
							SequenceNumber: int32(sqNumber),
							Timestamp:      time.Now().UnixMilli(),
							Format:         protobuf.VideoFormat_MJPEG,
							IsLastChunk:    false,
						},
					},
				}
				success := ss.PresenceSystem.TransmitionService.SendPacket(header, ss.PresenceSystem.RoutingTable)

				if !success {
					ss.Logger.Error(fmt.Sprintf("Failed to send video chunk to %s", header.GetSender()))
				}

				// Sleep to control frame rate (e.g., 30 FPS)
				time.Sleep(time.Millisecond * 33)
			}
		}

		// Properly stop ffmpeg if we break out of the loop
		cmd.Process.Kill()
		cmd.Wait()
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
