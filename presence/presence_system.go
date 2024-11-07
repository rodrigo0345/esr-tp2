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
	RoutingTable     *dvr.DistanceVectorRouting
	NeighborList     *NeighborList
	ConnectedClients *ClientList
	ConnectionPool   *dvr.ConnectionPool
	Logger           *config.Logger
	Config           *config.AppConfigList
}

func NewPresenceSystem(cnf *config.AppConfigList) *PresenceSystem {

	neighborList := &NeighborList{
		Content: cnf.Neighbors,
	}
	routingTable := dvr.CreateDistanceVectorRouting(cnf)
	connectedClients := &ClientList{
		Content: []*protobuf.Interface{},
	}

	// used to reduce the time spent opening and closing connections
	neighborsConnectionsMap := dvr.NewNeighborsConnectionsMap()

	logger := config.NewLogger(2)

	return &PresenceSystem{
		RoutingTable:     routingTable,
		NeighborList:     neighborList,
		ConnectedClients: connectedClients,
		ConnectionPool:   neighborsConnectionsMap,
		Logger:           logger,
		Config:           cnf,
	}
}

func (ps *PresenceSystem) HeartBeatNeighbors(seconds int) {
	for {
		ps.RoutingTable = ps.NeighborList.PingNeighbors(ps.Config, ps.RoutingTable, ps.ConnectionPool)
		ps.RoutingTable.Print()
		time.Sleep(time.Second * time.Duration(seconds))
	}
}

// listen for clients with QUIC
func (ps *PresenceSystem) ListenForClients() {
	tlsConfig := config.GenerateTLS()

	listener, err := quic.ListenAddr(dvr.Interface{Interface: ps.Config.NodeIP}.ToString(), tlsConfig, nil)

	if err != nil {
		ps.Logger.Error(err.Error())
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

			case protobuf.RequestType_RETRANSMIT:

				// check if the target is the client and check if the client is connected on this node
				if ps.ConnectedClients.Has(header.GetClientIp()) {

					// send the message to the client
					streaming.RedirectPacketToClient(data, header.GetClientIp())

					// remove the client from the list
					ps.ConnectedClients.Remove(header.GetClientIp())
					return
				}

				HandleRetransmit(ps, connection, stream, header)
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

			// add the client to the list
			ps.ConnectedClients.mutex.Lock()
			ps.ConnectedClients.Content = append(ps.ConnectedClients.Content, config.ToInterface(header.ClientIp))
			ps.ConnectedClients.mutex.Unlock()

      ps.Logger.Info(fmt.Sprintf("Received message from %s\n", addr))

			// Process the message based on its type
			switch header.Type {
			case protobuf.RequestType_RETRANSMIT:
				// modify data to set the sender to this node
				header.Sender = ps.Config.NodeName

				HandleRetransmitFromClient(ps, header)
			default:
        ps.Logger.Error(fmt.Sprintf("Received unrecognized packet type from %v\n", addr))
			}
		}(buffer[:n], remoteAddr)
	}


}
