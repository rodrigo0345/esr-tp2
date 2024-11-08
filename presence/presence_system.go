package presence

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	dvr "github.com/rodrigo0345/esr-tp2/presence/distance_vector_routing"
	"github.com/rodrigo0345/esr-tp2/presence/streaming"
)

type PresenceSystem struct {
	RoutingTable   *dvr.DistanceVectorRouting
	NeighborList   *NeighborList
	CurrentStreams map[string]*ClientList
	ConnectionPool *dvr.ConnectionPool
	Logger         *config.Logger
	Config         *config.AppConfigList
}

func NewPresenceSystem(cnf *config.AppConfigList) *PresenceSystem {
	neighborList := &NeighborList{
		Content: cnf.Neighbors,
	}
	routingTable := dvr.CreateDistanceVectorRouting(cnf)

	// used to reduce the time spent opening and closing connections
	neighborsConnectionsMap := dvr.NewNeighborsConnectionsMap()

	logger := config.NewLogger(2)

	return &PresenceSystem{
		RoutingTable:   routingTable,
		NeighborList:   neighborList,
		CurrentStreams: make(map[string]*ClientList),
		ConnectionPool: neighborsConnectionsMap,
		Logger:         logger,
		Config:         cnf,
	}
}

func (ps *PresenceSystem) HeartBeatNeighbors(seconds int) {
	for {
		ps.RoutingTable = ps.NeighborList.PingNeighbors(ps.Config, ps.RoutingTable, ps.ConnectionPool)
		ps.RoutingTable.Print()
		time.Sleep(time.Second * time.Duration(seconds))
	}
}

func (ps *PresenceSystem) HeartBeatClients(seconds int) {
	for {
		// kill clients that dont ping in a while
		for video := range ps.CurrentStreams {
			ps.CurrentStreams[video].RemoveDeadClients(seconds)

			// if the list is empty, notify the server that it can stop streaming
			if len(ps.CurrentStreams[video].Content) == 0 {
				ps.Logger.Info(fmt.Sprintf("No clients are connected to %s\n", video))
			}

			delete(ps.CurrentStreams, video)

			header := &protobuf.Header{
				Type:           protobuf.RequestType_RETRANSMIT,
				Length:         0,
				Timestamp:      int32(time.Now().UnixMilli()),
				ClientIp:       ps.Config.NodeIP.String(),
				Sender:         ps.Config.NodeName,
				Target:         "server", // TODO: change this and make it dynamic
				RequestedVideo: video,
				Content: &protobuf.Header_ClientCommand{
					ClientCommand: &protobuf.ClientCommand{
						Command:               protobuf.PlayerCommand_STOP,
						AdditionalInformation: "stopped because no clients are connected",
					},
				},
			}

			SendMessage(ps, header)
		}
		time.Sleep(time.Second * time.Duration(seconds))
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

			switch header.Type {
			case protobuf.RequestType_ROUTINGTABLE:

				HandleRouting(ps, connection, stream, header)
				break

			case protobuf.RequestType_RETRANSMIT:

				videoName := header.RequestedVideo

				// check if the target is the client and check if the client is connected on this node
				if ps.CurrentStreams[videoName] != nil && len(ps.CurrentStreams[videoName].Content) > 0 {

					for _, client := range ps.CurrentStreams[videoName].Content {
						// send the message to all interested clients
						streaming.RedirectPacketToClient(data, fmt.Sprintf("%s:%d", client.Ip, client.Port))
					}

					return
				}

				HandleRetransmit(ps, header)
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

		go func(data []byte, addr *net.UDPAddr) {
			header, err := config.UnmarshalHeader(data)
			if err != nil {
				ps.Logger.Error(err.Error())
				return
			}

			// Process the message based on its type
			switch header.Type {
			case protobuf.RequestType_RETRANSMIT:
				videoName := header.RequestedVideo

				// add the client to the list
				ps.CurrentStreams[videoName] = &ClientList{}
				ps.CurrentStreams[videoName].Add(config.ToInterface(header.ClientIp))
				ps.Logger.Info(fmt.Sprintf("Received message from %s\n", addr))

				// modify data to set the sender to this node
				header.Sender = ps.Config.NodeName

				HandleRetransmitFromClient(ps, header)

				// listen for clients in UDP, if their message takes more than 10 seconds, remove them from the NeighborList
				// and if the list is empty, notify the server that it can stop streaming
			case protobuf.RequestType_HEARTBEAT:

				for _, video := range ps.CurrentStreams {
					video.ResetPing(header.ClientIp)
				}

			default:
				ps.Logger.Error(fmt.Sprintf("Received unrecognized packet type from %v\n", addr))
			}
		}(buffer[:n], remoteAddr)
	}
}
