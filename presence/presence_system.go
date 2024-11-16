package presence

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	dvr "github.com/rodrigo0345/esr-tp2/presence/distance_vector_routing"
	"github.com/rodrigo0345/esr-tp2/presence/streaming"
)

type PresenceSystem struct {
	RoutingTable         *dvr.DistanceVectorRouting
	NeighborList         *NeighborList
	CurrentClientStreams map[string]*ClientList
	CurrentNodeStreams   map[string]*NodeList
	ConnectionPool       *dvr.ConnectionPool
	Logger               *config.Logger
	Config               *config.AppConfigList
}

func NewPresenceSystem(cnf *config.AppConfigList) *PresenceSystem {
	neighborList := &NeighborList{
		Content: cnf.Neighbors,
	}
	routingTable := dvr.CreateDistanceVectorRouting(cnf)

	// used to reduce the time spent opening and closing connections
	neighborsConnectionsMap := dvr.NewNeighborsConnectionsMap()
	currentNodeStreams := make(map[string]*NodeList)

	logger := config.NewLogger(2)

	return &PresenceSystem{
		RoutingTable:         routingTable,
		NeighborList:         neighborList,
		CurrentClientStreams: make(map[string]*ClientList),
		CurrentNodeStreams:   currentNodeStreams,
		ConnectionPool:       neighborsConnectionsMap,
		Logger:               logger,
		Config:               cnf,
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
		for video := range ps.CurrentClientStreams {
			ps.CurrentClientStreams[video].RemoveDeadClients(seconds)

			// if the list is empty, notify the server that it can stop streaming
			if len(ps.CurrentClientStreams[video].Content) != 0 {
				break
			}

			ps.Logger.Info(fmt.Sprintf("No clients are connected to %s\n", video))
			delete(ps.CurrentClientStreams, video)

			header := &protobuf.Header{
				Type:           protobuf.RequestType_RETRANSMIT,
				Length:         0,
				Timestamp:      time.Now().UnixMilli(),
				ClientIp:       fmt.Sprintf("%s:%d", ps.Config.NodeIP.Ip, ps.Config.NodeIP.Port),
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
			header.Length = int32(proto.Size(header))

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

        ps.Logger.Info(fmt.Sprintf("Received message from %s", header.Sender))
				videoName := header.RequestedVideo

				serverMessage := header.GetServerVideoChunk() != nil
				clientMessage := header.GetClientCommand() != nil

				// check if the target is the client and check if the client is connected on this node
				if serverMessage && ps.CurrentClientStreams[videoName] != nil && len(ps.CurrentClientStreams[videoName].Content) > 0 {

					for _, client := range ps.CurrentClientStreams[videoName].Content {
						ps.Logger.Info(fmt.Sprintf("Sending video chunk to %s\n", fmt.Sprintf("%s:%d", client.Ip, client.Port)))
						// send the message to all interested clients
						streaming.RedirectPacketToClient(data, fmt.Sprintf("%s:%d", client.Ip, client.Port))
					}

					break
				}

				// ------- EFFICIENT RETRANSMITION -------

				// make the nodes know what video they are about to stream
				if clientMessage && header.RequestedVideo != "" {
					ps.Logger.Info(fmt.Sprintf("Client %s is requesting video '%s'\n", header.GetSender(), header.RequestedVideo))
					// check if the node is already streaming and just add the client to the list
					if exists := ps.CurrentNodeStreams[header.RequestedVideo]; exists != nil {
						// now the current stream is going to be streamd to this node too

						ps.CurrentNodeStreams[header.RequestedVideo] = &NodeList{}
						ps.CurrentNodeStreams[header.RequestedVideo].Add(header.Sender)
					} else {
						ps.CurrentNodeStreams[header.RequestedVideo] = &NodeList{}
						ps.CurrentNodeStreams[header.RequestedVideo].Add(header.Sender)
					}
				}

				// this is the stream
				if serverMessage {

					// streaming
					video := header.RequestedVideo

					ps.Logger.Info(fmt.Sprintf("Server %s is sending video '%s'\n", header.GetSender(), header.RequestedVideo))
					if ps.CurrentNodeStreams[video] == nil {
						break
					}

					for _, node := range ps.CurrentNodeStreams[video].Content {
						// send the message to all interested clients
						ps.Logger.Info(fmt.Sprintf("Sending video chunk to %s\n", node.Node))
						header.Target = node.Node
						HandleRetransmit(ps, header)
					}

					break
				}

				// ------- EFFICIENT RETRANSMITION -------

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

		remoteIp := fmt.Sprintf("%s:%d", remoteAddr.IP.String(), 2222)

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
				ps.CurrentClientStreams[videoName] = &ClientList{}
				ps.CurrentClientStreams[videoName].Add(config.ToInterface(remoteIp))
				ps.Logger.Info(fmt.Sprintf("Received message from %s\n", addr))

				// modify data to set the sender to this node
				header.Sender = ps.Config.NodeName

				HandleRetransmitFromClient(ps, header)

				// listen for clients in UDP, if their message takes more than 10 seconds, remove them from the NeighborList
				// and if the list is empty, notify the server that it can stop streaming
			case protobuf.RequestType_HEARTBEAT:

				for _, video := range ps.CurrentClientStreams {
					video.ResetPing(header.ClientIp)
				}

			default:
				ps.Logger.Error(fmt.Sprintf("Received unrecognized packet type from %v\n", addr))
			}
		}(buffer[:n], remoteAddr)
	}
}
