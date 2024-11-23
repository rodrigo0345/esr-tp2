package presence

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"github.com/rodrigo0345/esr-tp2/presence/clientStreaming"
	dvr "github.com/rodrigo0345/esr-tp2/presence/distance_vector_routing"
	"github.com/rodrigo0345/esr-tp2/presence/transmitions"
)

type PresenceSystem struct {
	RoutingTable       *dvr.DistanceVectorRouting
	NeighborList       *NeighborList
	ConnectionPool     *dvr.ConnectionPool
	Logger             *config.Logger
	Config             *config.AppConfigList
	ClientService      *clientStreaming.StreamingService
	TransmitionService *transmitions.TransmissionService
}

func NewPresenceSystem(cnf *config.AppConfigList) *PresenceSystem {
	neighborList := &NeighborList{
		Content: cnf.Neighbors,
	}
	routingTable := dvr.CreateDistanceVectorRouting(cnf)

	// used to reduce the time spent opening and closing connections
	neighborsConnectionsMap := dvr.NewNeighborsConnectionsMap()

	logger := config.NewLogger(2, cnf.NodeName)

	streamingService := clientStreaming.NewStreamingService(cnf, logger)
	streamingService.RunBackgroundRoutine(5)

	transmissionService := transmitions.NewTransmissionService(logger, cnf)

	return &PresenceSystem{
		RoutingTable:       routingTable,
		NeighborList:       neighborList,
		ConnectionPool:     neighborsConnectionsMap,
		Logger:             logger,
		Config:             cnf,
		ClientService:      streamingService,
		TransmitionService: transmissionService,
	}
}

func (ps *PresenceSystem) HeartBeatNeighbors(seconds int) {
	for {
		ps.RoutingTable = ps.NeighborList.PingNeighbors(ps.Logger, ps.Config, ps.RoutingTable, ps.ConnectionPool)
		// ps.RoutingTable.Print(ps.Logger)
		time.Sleep(time.Second * time.Duration(seconds))
	}
}

func (ps *PresenceSystem) SignalDeadClientsService() {
	for {
		select {
		case data := <-ps.ClientService.SignalDead:
			ps.Logger.Info(fmt.Sprintf("Client %s is not pinging\n", data.Header.GetSender()))
			ps.TransmitionService.SendPacket(data.Header, ps.RoutingTable)
		}
	}
}

// listen for clients with QUIC
func (ps *PresenceSystem) ListenForClients() {
	tlsConfig := config.GenerateTLS()

	listener, err := quic.ListenAddr(dvr.Interface{Interface: ps.Config.NodeIP}.ToString(), tlsConfig, nil)

	if err != nil {
		ps.Logger.Error(err.Error())
		return
	}

	ps.Logger.Info(fmt.Sprintf("QUIC server is listening on port %d\n", ps.Config.NodeIP.Port))

	for {
		connection, err := listener.Accept(context.Background())
		if err != nil {
			ps.Logger.Error(err.Error())
			continue
		}

		go func(session quic.Connection) {

			stream, err := session.AcceptStream(context.Background())
			if err != nil {
				ps.Logger.Error(err.Error())
				return
			}

			data, err := config.ReceiveMessage(stream)
			if err != nil {
				ps.Logger.Error(err.Error())
				return
			}

			header, err := config.UnmarshalHeader(data)
			if err != nil {
				ps.Logger.Error(err.Error())
				return
			}

			header.Path = fmt.Sprintf("%s,%s", header.Path, ps.Config.NodeName)

			switch header.Type {
			case protobuf.RequestType_ROUTINGTABLE:

				HandleRouting(ps, connection, stream, header)

				break

			case protobuf.RequestType_RETRANSMIT:

				isVideoPacket := header.GetServerVideoChunk() != nil

				if !isVideoPacket {
					// retransmit the packet
					ps.TransmitionService.SendPacket(header, ps.RoutingTable)
				}

				var callback chan clientStreaming.CallbackData = make(chan clientStreaming.CallbackData)
				ps.ClientService.Signal <- clientStreaming.SignalData{
					Command:  clientStreaming.VIDEO,
					Packet:   header,
					Callback: callback,
				}

				// it already sends the packet to the client
				select {
				case data := <-callback:
					// not for us so, send to another neighbor
					if data.Cancel {
						ps.Logger.Info(fmt.Sprintf("Forwarding request from %s", header.GetSender()))
						ps.TransmitionService.SendPacket(header, ps.RoutingTable)
					} else {
						ps.Logger.Info(fmt.Sprintf("Reached the client"))
					}
				}

				break
			default:

			}
		}(connection)
	}

}

func (ps *PresenceSystem) ListenForClientsInUDP() {
	addr := ps.Config.NodeIP
	udpPort := ps.Config.NodeIP.Port - 1000
	addrString := fmt.Sprintf("%s:%d", addr.Ip, udpPort)

	var netAddr *net.UDPAddr
	netAddr, err := net.ResolveUDPAddr("udp", addrString)
	conn, err := net.ListenUDP("udp", netAddr)

	if err != nil {
		ps.Logger.Error(err.Error())
	}
	defer conn.Close()

	ps.Logger.Info(fmt.Sprintf("UDP server is listening on %s, this is only used for direct connection with clients\n", addrString))

	buffer := make([]byte, 1024) // Buffer size to read incoming packets
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			ps.Logger.Error(err.Error())
			continue
		}

		remoteIp := fmt.Sprintf("%s:%d", remoteAddr.IP.String(), 2222)

		go func(data []byte, addr *net.UDPAddr) {
			header, err := config.UnmarshalHeader(data)
			if err != nil {
				ps.Logger.Error(err.Error())
				return
			}

			videoName := header.RequestedVideo

			// Process the message based on its type
			switch header.Type {
			case protobuf.RequestType_RETRANSMIT:

				isRequestingPlay := header.GetClientCommand().Command == protobuf.PlayerCommand_PLAY

				var operation clientStreaming.Command
				if isRequestingPlay {
					operation = clientStreaming.PLAY
				} else {
					operation = clientStreaming.STOP
				}

				ps.Logger.Info(fmt.Sprintf("Client %s is requesting video '%s'\n", header.GetSender(), header.RequestedVideo))

				var callback chan clientStreaming.CallbackData = make(chan clientStreaming.CallbackData)
				ps.ClientService.Signal <- clientStreaming.SignalData{
					Command:   operation,
					Video:     clientStreaming.Video(videoName),
					UdpClient: config.ToInterface(remoteIp),
					Callback:  callback,
				}

				select {
				case data := <-callback:
					if data.Cancel {

						ps.Logger.Info(fmt.Sprintf("Client %s is already connected to the server\n", addr))
						return

					} else {

						success := ps.TransmitionService.SendPacket(data.Header, ps.RoutingTable)

						if !success {
							// stop
							ps.ClientService.SendSignal(clientStreaming.SignalData{
								Command:   clientStreaming.STOP,
								Video:     clientStreaming.Video(videoName),
								UdpClient: config.ToInterface(remoteIp),
								Callback:  nil,
							})
						}

					}
				}

			case protobuf.RequestType_HEARTBEAT:

				ps.ClientService.SendSignal(clientStreaming.SignalData{
					Command:   clientStreaming.PING,
					Video:     clientStreaming.Video(videoName),
					UdpClient: config.ToInterface(remoteIp),
					Callback:  nil,
				})

				select {
				case data := <-ps.ClientService.SignalDead:
					// isto não faz nada
					ps.Logger.Info(fmt.Sprintf("Client %s is not pinging\n", data.Header.GetSender()))
				default:
				}

			default:
				ps.Logger.Error(fmt.Sprintf("Received unrecognized packet type from %v\n", addr))
			}
		}(buffer[:n], remoteAddr)
	}
}
