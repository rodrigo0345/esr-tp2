package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"github.com/rodrigo0345/esr-tp2/presence"
	"github.com/rodrigo0345/esr-tp2/server/boostrapper"
	"google.golang.org/protobuf/proto"
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

func (ss *ServerSystem) BsClients() {
	addr := ss.PresenceSystem.Config.NodeIP
	udpPort := addr.Port - 1000 // Assuming NodeIP.Port is an integer
	addrString := fmt.Sprintf("%s:%d", addr.Ip, udpPort)

	netAddr, err := net.ResolveUDPAddr("udp", addrString)
	if err != nil {
		ss.Logger.Error(fmt.Sprintf("Failed to resolve UDP address %s: %v", addrString, err))
		return
	}

	conn, err := net.ListenUDP("udp", netAddr)
	if err != nil {
		ss.Logger.Error(fmt.Sprintf("Failed to listen on UDP address %s: %v", addrString, err))
		return
	}
	defer conn.Close()
	ss.Logger.Info(fmt.Sprintf("Listening for UDP retransmit on %s", addrString))

	for {
		data, remoteAddr, err := config.ReceiveMessageUDP(conn)
		if err != nil {
			ss.Logger.Error(fmt.Sprintf("Error reading from UDP: %v", err))
			continue // Consider breaking the loop or implementing a retry mechanism based on the error
		}

		remoteIp := fmt.Sprintf("%s:%d", remoteAddr.IP, 2222)

		go func(data []byte, remoteAddr string) {

			header := &protobuf.Header{}
			err = proto.Unmarshal(data, header)

			if err != nil {
				ss.Logger.Error(fmt.Sprintf("Error unmarshalling UDP message: %v", err))
				return
			}

			switch header.Type {
			case protobuf.RequestType_CLIENT_PING:
				ss.Bootstrapper.BootstrapClients(header, remoteIp)
				break
			default:
			}
		}(data, remoteIp)

	}
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
		// Initialize ffmpeg command to stream video as MJPEG
		cmd := exec.Command("ffmpeg",
			"-stream_loop", "-1",
			"-i", videoFilePath,
			"-f", "image2pipe",
			"-vcodec", "mjpeg",
			"-",
		)

		// Capture stdout (video frames)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			ss.Logger.Error(fmt.Sprintf("Failed to get stdout pipe for video %s: %v", video.Video, err))
			return
		}

		// Capture stderr (metrics)
		stderr, err := cmd.StderrPipe()
		if err != nil {
			ss.Logger.Error(fmt.Sprintf("Failed to get stderr pipe for video %s: %v", video.Video, err))
			return
		}

		if err := cmd.Start(); err != nil {
			ss.Logger.Error(fmt.Sprintf("Failed to start ffmpeg for video %s: %v", video.Video, err))
			return
		}

		reader := bufio.NewReader(stdout)
		metricsReader := bufio.NewReader(stderr)

		// Use WaitGroup to manage goroutines
		var wg sync.WaitGroup
		metricsChan := make(chan ServerMetrics, 1)
		wg.Add(1)
		go func() {
			defer wg.Done()
			collectMetrics(metricsReader, metricsChan)
		}()

		sqNumber := 0

		// Initialize default metrics
		currentMetrics := ServerMetrics{
			BitrateKbps:         0,
			Width:               0,
			Height:              0,
			Fps:                 0.0,
			LatencyMs:           0,
			CurrentBitrateKbps:  0,
			PreviousBitrateKbps: 0,
		}

		for {
			select {
			case <-stopChan:
				cmd.Process.Kill() // Terminate ffmpeg process
				ss.Logger.Info(fmt.Sprintf("Stream stopped for video %s", video.Video))
				wg.Wait()
				return
			case metrics := <-metricsChan:
				// Update current metrics
				currentMetrics = metrics
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

				ss.Logger.Info(fmt.Sprintf("Sending video chunk to %s", targets))

				// Populate the ServerVideoChunk with dynamic metrics
				serverVideoChunk := &protobuf.ServerVideoChunk{
					SequenceNumber:      int32(sqNumber),
					Timestamp:           time.Now().UnixMilli(),
					Format:              protobuf.VideoFormat_MJPEG,
					Data:                frameData,
					IsLastChunk:         false,
					BitrateKbps:         currentMetrics.BitrateKbps,
					Width:               int32(currentMetrics.Width),
					Height:              int32(currentMetrics.Height),
					Fps:                 currentMetrics.Fps,
					LatencyMs:           currentMetrics.LatencyMs,
					CurrentBitrateKbps:  currentMetrics.CurrentBitrateKbps,
					PreviousBitrateKbps: currentMetrics.PreviousBitrateKbps,
				}

				header := &protobuf.Header{
					Sender:         ss.PresenceSystem.Config.NodeName,
					Target:         targets,
					RequestedVideo: video.Video,
					Path:           ss.PresenceSystem.Config.NodeName,
					Type:           protobuf.RequestType_RETRANSMIT,
					Length:         int32(len(frameData)),
					MaxHops:        10,
					Hops:           1,
					Timestamp:      time.Now().UnixMilli(),
					Content: &protobuf.Header_ServerVideoChunk{
						ServerVideoChunk: serverVideoChunk,
					},
				}

				success := ss.PresenceSystem.TransmitionService.SendPacket(header, ss.PresenceSystem.RoutingTable, true)

				if !success {
					ss.Logger.Error(fmt.Sprintf("Failed to send video chunk to %s", header.GetSender()))
				}

				time.Sleep(time.Millisecond * 33)
			}
		}

		// Properly stop ffmpeg if we break out of the loop
		cmd.Process.Kill()
		cmd.Wait()
		wg.Wait()
	}
}

// ServerMetrics holds the dynamic metrics collected from ffmpeg
type ServerMetrics struct {
	BitrateKbps         int32
	Width               int
	Height              int
	Fps                 float64
	LatencyMs           int64
	CurrentBitrateKbps  int32
	PreviousBitrateKbps int32
}

// collectMetrics parses ffmpeg's stderr output to extract metrics
func collectMetrics(reader *bufio.Reader, metricsChan chan<- ServerMetrics) {
	// Regular expressions to match ffmpeg's metric output
	// Example ffmpeg stderr output line: "frame=  240 fps=25 q=28.0 size=    1024kB time=00:00:10.00 bitrate= 838.9kbits/s speed=1.00x"
	frameRegex := regexp.MustCompile(`frame=\s*(\d+).*fps=([\d.]+).*size=\s*(\d+)kB.*bitrate=\s*([\d.]+)kbits/s`)
	timeRegex := regexp.MustCompile(`time=(\d+:\d+:\d+\.\d+)`)

	var (
		previousBitrate float64
	)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			// Log and continue on other errors
			fmt.Printf("Error reading ffmpeg stderr: %v\n", err)
			continue
		}

		// Parse frame metrics
		if matches := frameRegex.FindStringSubmatch(line); matches != nil {
			// frameNumber, _ := strconv.Atoi(matches[1])
			fps, _ := strconv.ParseFloat(matches[2], 64)
			// sizeKB, _ := strconv.Atoi(matches[3])
			bitrate, _ := strconv.ParseFloat(matches[4], 64)

			// Example calculation for latency (requires additional implementation)
			latencyMs := calculateLatency()

			metrics := ServerMetrics{
				BitrateKbps:         int32(bitrate),
				Fps:                 fps,
				LatencyMs:           latencyMs,
				CurrentBitrateKbps:  int32(bitrate),
				PreviousBitrateKbps: int32(previousBitrate),
			}

			// Update previous bitrate
			previousBitrate = bitrate

			metricsChan <- metrics
		}

		// Optionally, parse other metrics like time or resolution if available
		if matches := timeRegex.FindStringSubmatch(line); matches != nil {
			// Implement parsing if needed
		}
	}
}

// calculateLatency is a placeholder for latency calculation logic
func calculateLatency() int64 {
	// Implement latency calculation logic here
	// This could involve tracking timestamps when frames are sent and received
	return 0 // Placeholder value
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
