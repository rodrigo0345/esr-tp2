package presence

import (
	"context"
	"crypto/tls"
	_ "encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	distancevectorrouting "github.com/rodrigo0345/esr-tp2/presence/distance_vector_routing"
	"github.com/rodrigo0345/esr-tp2/presence/streaming"
	"github.com/rodrigo0345/esr-tp2/server/certs"
	"google.golang.org/protobuf/proto"
)

type NeighborList struct {
	mutex   sync.Mutex
	Content []*protobuf.Interface
}

func (nl *NeighborList) Has(neighbor *protobuf.Interface) bool {
	nl.mutex.Lock()
	defer nl.mutex.Unlock()

	for _, n := range nl.Content {
		if n.String() == neighbor.String() {
			return true
		}
	}

	return false
}

type ClientList struct {
	mutex   sync.Mutex
	Content []*protobuf.Interface
}

func (cl *ClientList) Remove(client string) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	for i, c := range cl.Content {
		if c.String() == client {
			cl.Content = append(cl.Content[:i], cl.Content[i+1:]...)
			break
		}
	}
}

func (cl *ClientList) Add(client *protobuf.Interface) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	cl.Content = append(cl.Content, client)
}

func (cl *ClientList) Has(clientIp string) bool {
	clIP := config.ToInterface(clientIp)
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	for _, c := range cl.Content {
		if c.Ip == clIP.Ip && c.Port == clIP.Port {
			return true
		}
	}
	return false
}

func Presence(config *config.AppConfigList) {
	if config.NodeIP == nil {
		log.Fatal("ServerUrl is nil")
	}

	neighborList := &NeighborList{
		Content: config.Neighbors,
	}
	routingTable := distancevectorrouting.CreateDistanceVectorRouting(config)
	connectedClients := &ClientList{
		Content: []*protobuf.Interface{},
	}

	// used to reduce the time spent opening and closing connections
	neighborsConnectionsMap := distancevectorrouting.NewNeighborsConnectionsMap()

	go func() {
		for {
			routingTable = neighborList.PingNeighbors(config, routingTable, neighborsConnectionsMap)
			// routingTable.Print()
			time.Sleep(time.Second * 6)
		}
	}()

	fmt.Printf("Node is running on %s\n", distancevectorrouting.Interface{Interface: config.NodeIP}.ToString())

	go MainListen(config, neighborList, neighborsConnectionsMap, routingTable, connectedClients)
	ListenForClientsInUDP(config, neighborList, neighborsConnectionsMap, routingTable, connectedClients)
}

func MainListen(cnf *config.AppConfigList, neighborList *NeighborList, neighborsConnectionsMap *distancevectorrouting.ConnectionPool, routingTable *distancevectorrouting.DistanceVectorRouting, connectedClients *ClientList) {
	tlsConfig := generateTLS()

	// Wait for connections
	listener, err := quic.ListenAddr(distancevectorrouting.Interface{Interface: cnf.NodeIP}.ToString(), tlsConfig, nil)
	if err != nil {
		log.Fatalf("Failed to start QUIC listener: %v", err)
	}

	fmt.Printf("QUIC server is listening on port %d\n", cnf.NodeIP.Port)

	for {
		connection, err := listener.Accept(context.Background())
		if err != nil {
			log.Printf("Failed to accept session: %v", err)
			continue
		}

		go func(session quic.Connection) {

			stream, err := session.AcceptStream(context.Background())
			if err != nil {
				log.Printf("Failed to accept stream for session %v: %v", session.RemoteAddr(), err)
				return
			}

			data, err := config.ReceiveMessage(stream)
			if err != nil {
				log.Printf("Failed to receive message from %v: %v", stream.StreamID(), err)
				return
			}

			chunk := &protobuf.Header{}
			err = proto.Unmarshal(data, chunk)
			if err != nil {
				log.Printf("Failed to unmarshal message from %v: %v", stream.StreamID(), err)
				return
			}

			switch chunk.Type {
			case protobuf.RequestType_ROUTINGTABLE:
				timeTook := int32(time.Since(time.UnixMilli(int64(chunk.Timestamp))))
				otherDvr := &distancevectorrouting.DistanceVectorRouting{Mutex: sync.Mutex{}, Dvr: chunk.GetDistanceVectorRouting()}
				HandleRouting(session, cnf, stream, neighborList, routingTable, otherDvr, timeTook)
			case protobuf.RequestType_RETRANSMIT:

				// check if the target is the client and check if the client is connected on this node
				if connectedClients.Has(chunk.GetClientIp()) {

					// send the message to the client
					streaming.RedirectPacketToClient(data, chunk.GetClientIp())

					// remove the client from the list
					connectedClients.Remove(chunk.GetClientIp())

					return
				}

				HandleRetransmit(neighborList, routingTable, chunk, neighborsConnectionsMap)
			}
		}(connection)
	}
}

func ListenForClientsInUDP(cnf *config.AppConfigList, neighborList *NeighborList, neighborsConnectionsMap *distancevectorrouting.ConnectionPool, routingTable *distancevectorrouting.DistanceVectorRouting, connectedClients *ClientList) {
	addr := cnf.NodeIP
	udpPort := cnf.NodeIP.Port - 1000
	addrString := fmt.Sprintf("%s:%d", addr.Ip, udpPort)

	var netAddr *net.UDPAddr
	netAddr, err := net.ResolveUDPAddr("udp", addrString)
	conn, err := net.ListenUDP("udp", netAddr)

	if err != nil {
		log.Fatalf("Failed to start UDP listener: %v", err)
	}
	defer conn.Close()

	fmt.Printf("UDP server is listening on %s, this is only used for direct connection with clients\n", addrString)

	buffer := make([]byte, 1024) // Buffer size to read incoming packets
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Failed to read from UDP: %v", err)
			continue
		}

		go func(data []byte, addr *net.UDPAddr) {
			chunk := &protobuf.Header{}
			err := proto.Unmarshal(data[:n], chunk)
			if err != nil {
				log.Printf("Failed to unmarshal message from %v: %v", addr, err)
				return
			}

			// add the client to the list
			connectedClients.mutex.Lock()
			connectedClients.Content = append(connectedClients.Content, config.ToInterface(chunk.ClientIp))
			connectedClients.mutex.Unlock()

			fmt.Printf("Received message from %v: %v\n", addr, string(data[:n]))

			// Process the message based on its type
			switch chunk.Type {
			case protobuf.RequestType_RETRANSMIT:
				// modify data to set the sender to this node
				chunk.Sender = cnf.NodeName

				HandleRetransmitFromClient(cnf, neighborList, routingTable, chunk, neighborsConnectionsMap)
			default:
				log.Printf("Received unrecognized packet type from %v", addr)
			}
		}(buffer[:n], remoteAddr)
	}

}
func generateTLS() *tls.Config {

	certFile, keyFile, err := certs.GenerateTLSFiles()

	if err != nil {
		log.Fatalf("Failed to generate TLS files: %v", err)
	}
	defer os.Remove(certFile)
	defer os.Remove(keyFile)

	// Load the TLS configuration
	tlsConfig, err := certs.GenerateTLSConfig(certFile, keyFile)

	if err != nil {
		log.Fatalf("Failed to load TLS configuration: %v", err)
	}
	return tlsConfig
}
