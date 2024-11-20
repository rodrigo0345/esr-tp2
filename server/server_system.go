package server

import (
	"context"
	"fmt"
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
}

func NewServerSystem(cnf *config.AppConfigList) *ServerSystem {
	return &ServerSystem{
		PresenceSystem: presence.NewPresenceSystem(cnf),
		Logger:         config.NewLogger(2, cnf.NodeName),
    Bootstrapper:   boostrapper.NewBootstrapper("./server/boostrapper/nb.json", config.NewLogger(2, cnf.NodeName)),
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

				// unmarshal protobuf
				header, err := config.UnmarshalHeader(msg)

				if err != nil {
					ss.Logger.Error(err.Error())
					return
				}

				ss.Logger.Info(fmt.Sprintf("Received message from %s, type %s", header.GetSender(), header.GetType()))

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
						ss.Logger.Error(fmt.Sprintf("Received message with no client command\n"))
						return
					}

					if header.GetTarget() != ss.PresenceSystem.Config.NodeName {
						ss.Logger.Error(fmt.Sprintf("Received message from %s, but target is %s\n", header.GetSender(), header.GetTarget()))
						return
					}
					ss.Logger.Info(fmt.Sprintf("Received message %s\n", header.GetClientCommand().AdditionalInformation))

					// this specifies where the message needs to go
					header.Target = header.Sender
					header.Sender = ss.PresenceSystem.Config.NodeName
					stream.Close()

					switch header.GetClientCommand().Command {
					case protobuf.PlayerCommand_PLAY:
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
          break;
				}

			}(stream)
		}()
	}
}

func (ss *ServerSystem) StartVideoStream(header *protobuf.Header) {
	videoChoice := header.RequestedVideo
	clientName := header.GetTarget()
	ss.Logger.Info(fmt.Sprintf("Starting video stream for %s, requesting %s\n", clientName, videoChoice))

	client := &Client{PresenceNodeName: clientName, ClientIP: header.ClientIp}

	ss.VideoStreams.AddStream(videoChoice, client)
}

func (ss *ServerSystem) StopVideoStream(header *protobuf.Header) {
	videoChoice := header.RequestedVideo
	clientName := header.GetTarget()
	ss.Logger.Info(fmt.Sprintf("Stopping video stream for %s, requesting %s\n", clientName, videoChoice))

	client := Client{PresenceNodeName: clientName, ClientIP: header.ClientIp}
	ss.VideoStreams.RemoveStream(videoChoice, client)
}

func (ss *ServerSystem) BackgroundStreaming() {
	for {
		for _, video := range ss.VideoStreams.Streams {

			header := &protobuf.Header{
				Sender: ss.PresenceSystem.Config.NodeName,
				// to be defined
				Target:         "",
				RequestedVideo: video.Video,
				Type:           protobuf.RequestType_RETRANSMIT,
				// to be determined
				Length:    0,
				Timestamp: time.Now().UnixMilli(),
				Content: &protobuf.Header_ServerVideoChunk{
					ServerVideoChunk: &protobuf.ServerVideoChunk{
						Data:           []byte(fmt.Sprintf("Sending video stream for %s", video.Video)),
						SequenceNumber: 0,
						Timestamp:      time.Now().UnixMilli(),
						Format:         protobuf.VideoFormat_MJPEG,
						IsLastChunk:    false,
					},
				},
			}

			// TODO: stream video, for now just send a bunch of messages
			for _, client := range video.Clients {

				sqNumber := 0
				for len(video.Clients) != 0 {

					sqNumber += 1
					header.Target = client.PresenceNodeName
					header.ClientIp = client.ClientIP
					header.Length = int32(len(header.Content.(*protobuf.Header_ServerVideoChunk).ServerVideoChunk.Data))
					header.GetServerVideoChunk().SequenceNumber = int32(sqNumber)

					ss.sendVideoChunk(header)
				}
			}

		}
	}
}

func (ss *ServerSystem) sendVideoChunk(header *protobuf.Header) {
	target := header.GetTarget()

	nextHop, err := ss.PresenceSystem.RoutingTable.GetNextHop(target)
	if err != nil {
		ss.Logger.Error(err.Error())
		return
	}

	// open a connection with nextHop and be persistent trying to send the message
	neighbor := nextHop.NextNode

	failCount := 0
	limitFails := 3

	for {
		var msg []byte
		neighborIp := fmt.Sprintf("%s:%d", neighbor.Ip, neighbor.Port)
		neighborStream, _, err := ss.PresenceSystem.ConnectionPool.GetConnectionStream(neighborIp)
		defer config.CloseStream(neighborStream)

		if err != nil {
			ss.Logger.Error(err.Error())
			goto fail
		}

		msg, err = config.MarshalHeader(header)
		err = config.SendMessage(neighborStream, msg)

		if err != nil {
			ss.Logger.Error(err.Error())
			goto fail
		}

		break

	fail:
		failCount += 1
		if failCount > limitFails {
			ss.Logger.Error(fmt.Sprintf("Failed to send video chunk to %s\n", neighborIp))
			break
		}
		continue
	}
}
