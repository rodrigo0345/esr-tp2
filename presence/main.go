package presence

import (
	"context"
	"crypto/tls"
	_ "encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	distancevectorrouting "github.com/rodrigo0345/esr-tp2/presence/distance_vector_routing"
	"github.com/rodrigo0345/esr-tp2/server/certs"
	"google.golang.org/protobuf/proto"
)

type NeighborList struct {
	mutex   sync.Mutex
	content []*protobuf.Interface
}

func Presence(config *config.AppConfigList) {
	// create its table
	if config.NodeIP == nil {
		log.Fatal("ServerUrl is nil")
	}

	neighborList := &NeighborList{
		content: config.Neighbors,
	}
	routingTable := distancevectorrouting.CreateDistanceVectorRouting(distancevectorrouting.Interface{Interface: config.NodeIP})
	routingTable.Print()

	go func() {
		for {
			// send our routing table to all neighbors
			neighborList.PingNeighbors(config, routingTable)
			routingTable.Print()
			time.Sleep(time.Second * 10)
		}
	}()
	// identify and alert neighbors, if some neighbors are not reachable, they will be removed from the routing table

	fmt.Printf("Node is running on %s\n", distancevectorrouting.Interface{Interface: config.NodeIP}.ToString())

	MainListen(config, neighborList, routingTable)
}

func MainListen(cnf *config.AppConfigList, neighborList *NeighborList, routingTable *distancevectorrouting.DistanceVectorRouting) {
	fmt.Println("Listening")

	tlsConfig := generateTLS()

	// Wait for connections
	listener, err := quic.ListenAddr(distancevectorrouting.Interface{Interface: cnf.NodeIP}.ToString(), tlsConfig, nil)
	if err != nil {
		log.Fatalf("Failed to start QUIC listener: %v", err)
	}

	fmt.Printf("QUIC server is listening on port %d\n", cnf.NodeIP.Port)

	for {
		// Accept incoming QUIC connections
		session, err := listener.Accept(context.Background())
		if err != nil {
			log.Printf("Failed to accept session: %v", err)
			continue // Log error and continue accepting other sessions
		}

		// Handle the session in a goroutine
		go func(session quic.Connection) {
			defer session.CloseWithError(0, "session closed")

			// Accept only one stream for this session
			stream, err := session.AcceptStream(context.Background())
			if err != nil {
				log.Printf("Failed to accept stream for session %v: %v", session.RemoteAddr(), err)
				return
			}
			defer stream.Close()

			data, err := config.ReceiveMessage(stream)
			if err != nil {
				log.Printf("Failed to receive message from %v: %v", stream.StreamID(), err)
				return
			}

			// Unmarshal the protobuf message
			chunk := &protobuf.Header{}
			err = proto.Unmarshal(data, chunk)
			if err != nil {
				log.Printf("Failed to unmarshal message from %v: %v", stream.StreamID(), err)
				return
			}

			switch chunk.Type {
			case protobuf.RequestType_ROUTINGTABLE:
				otherDvr := &distancevectorrouting.DistanceVectorRouting{chunk.GetDistanceVectorRouting()}
				HandleRouting(stream, neighborList, routingTable, otherDvr, chunk.Timestamp)
			case protobuf.RequestType_RETRANSMIT:
				// TODO: handle retransmit
			}

			// chunkJSON, err := json.MarshalIndent(chunk, "", "  ")
			// if err != nil {
			// 	log.Printf("Failed to marshal chunk from %v: %v", session.RemoteAddr(), err)
			// 	return
			// }
			// log.Printf("Processed chunk from %v: %s", session.RemoteAddr(), chunkJSON)

			if err != nil {
				log.Printf("Failed to write message to %v: %v", stream.StreamID(), err)
				return
			}
		}(session)
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
