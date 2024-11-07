package server

import (
	"context"
	"fmt"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"github.com/rodrigo0345/esr-tp2/presence"
)

type ServerSystem struct {
	PresenceSystem *presence.PresenceSystem
	Logger         *config.Logger
}

func NewServerSystem(cnf *config.AppConfigList) *ServerSystem {
	return &ServerSystem{
		PresenceSystem: presence.NewPresenceSystem(cnf),
		Logger:         config.NewLogger(2),
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
				case protobuf.RequestType_RETRANSMIT:

					ss.Logger.Info(fmt.Sprintf("Received message %s\n", header.GetClientCommand().AdditionalInformation))
					header.Target = header.Sender
					header.Sender = ss.PresenceSystem.Config.NodeName

					header.Content = &protobuf.Header_ServerVideoChunk{
						ServerVideoChunk: &protobuf.ServerVideoChunk{
							Data: []byte(header.GetClientCommand().AdditionalInformation + "Hi there"),
						},
					}
					stream.Close()

					presence.HandleRetransmit(ss.PresenceSystem, session, stream, header)
				}

			}(stream)
		}()
	}
}
