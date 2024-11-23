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
	streamingService.RunBackgroundRoutine(120)

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

	// Start the QUIC listener
	listener, err := quic.ListenAddr(dvr.Interface{Interface: ps.Config.NodeIP}.ToString(), tlsConfig, nil)
	if err != nil {
		ps.Logger.Error(err.Error())
		return
	}

	ps.Logger.Info(fmt.Sprintf("QUIC server is listening on port %d", ps.Config.NodeIP.Port))

	for {
		// Accept incoming connections
		connection, err := listener.Accept(context.Background())
		if err != nil {
			ps.Logger.Error(fmt.Sprintf("Failed to accept connection: %v", err))
			continue
		}

		// Handle each connection in a separate goroutine
		go func(session quic.Connection) {

			// Accept a new stream from the connection
			stream, err := session.AcceptStream(context.Background())
			if err != nil {
				return
			}

			// Handle the stream in a goroutine
			go func(stream quic.Stream) {

				// Receive message data
				data, err := config.ReceiveMessage(stream)
				if err != nil {
					// Ignore the error, discard the stream
					ps.Logger.Debug(fmt.Sprintf("Failed to receive message: %v", err))
					return
				}

				// Unmarshal the message header
				header, err := config.UnmarshalHeader(data)
				if err != nil {
					ps.Logger.Debug(fmt.Sprintf("Failed to unmarshal header: %v", err))
					return
				}

				// Append the current node to the path
				header.Path = fmt.Sprintf("%s,%s", header.Path, ps.Config.NodeName)

				// Handle specific header types
				switch header.Type {
				case protobuf.RequestType_ROUTINGTABLE:
					HandleRouting(ps, connection, stream, header)

				case protobuf.RequestType_RETRANSMIT:

          config.CloseStream(stream)

					isVideoPacket := header.GetServerVideoChunk() != nil
					if !isVideoPacket {
						// Retransmit the packet to neighbors
						ps.TransmitionService.SendPacket(header, ps.RoutingTable)
					}

					// Send to clients via signal channel
					callback := make(chan clientStreaming.CallbackData, 1) // Buffered channel to avoid blocking
					ps.ClientService.Signal <- clientStreaming.SignalData{
						Command:  clientStreaming.VIDEO,
						Packet:   header,
						Callback: callback,
					}

					select {
					case data := <-callback:
						// If canceled, retransmit to neighbors
						if data.Cancel {
							ps.TransmitionService.SendPacket(header, ps.RoutingTable)
						}
					case <-time.After(time.Millisecond * 100): // Timeout to avoid indefinite blocking
						ps.Logger.Debug("Callback timed out, ignoring")
					}

				default:
					ps.Logger.Debug(fmt.Sprintf("Unknown packet type: %v", header.Type))
				}
			}(stream)
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
