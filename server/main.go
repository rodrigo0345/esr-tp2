package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"github.com/rodrigo0345/esr-tp2/presence"
	distancevectorrouting "github.com/rodrigo0345/esr-tp2/presence/distance_vector_routing"
	"github.com/rodrigo0345/esr-tp2/server/certs"
	"google.golang.org/protobuf/proto"
)

func Server(config *config.AppConfigList) {
	tls := generateTLS()
	config.NodeIP.Port = 1111

	// ping neighbours
	neighborList := &presence.NeighborList{
		Content: config.Neighbors,
	}
	routingTable := distancevectorrouting.CreateDistanceVectorRouting(config)

	// used to reduce the time spent opening and closing connections
	neighborsConnectionsMap := distancevectorrouting.NewNeighborsConnectionsMap()

	go func() {
		for {
			routingTable = neighborList.PingNeighbors(config, routingTable, neighborsConnectionsMap)
			routingTable.Print()
			time.Sleep(time.Second * 6)
		}
	}()

	listenAndServe(tls, int(config.NodeIP.Port), config, neighborList, routingTable)
}

func listenAndServe(tls *tls.Config, serverPort int, cnf *config.AppConfigList, nl *presence.NeighborList, rt *distancevectorrouting.DistanceVectorRouting) {
	port := fmt.Sprintf(":%d", serverPort)
	listener, err := quic.ListenAddr(port, tls, nil)
	if err != nil {
		log.Fatalf("Failed to start QUIC listener: %v", err)
	}
	fmt.Printf("QUIC server is listening on port %d\n", serverPort)

	for {
		// Accept incoming QUIC connections
		session, err := listener.Accept(context.Background())
		if err != nil {
			log.Fatalf("Failed to accept session: %v", err)
		}

		// For each session, handle incoming streams
		go func() {
			for {
				stream, err := session.AcceptStream(context.Background())
				if err != nil {
					log.Println("Failed to accept stream:", err)
					return
				}

				go func(stream quic.Stream) {
					msg, err := config.ReceiveMessage(stream)
					if err != nil {
						log.Println("Error receiving message:", err)
						return
					}

					// unmarshal protobuf
					var msgProtobuf protobuf.Header
					err = proto.Unmarshal(msg, &msgProtobuf)
					if err != nil {
						log.Println("Error unmarshal protobuf:", err)
						return
					}

					switch msgProtobuf.Type {
					case protobuf.RequestType_ROUTINGTABLE:
						timeTook := int32(time.Since(time.UnixMilli(int64(msgProtobuf.Timestamp))))
						otherDvr := &distancevectorrouting.DistanceVectorRouting{Mutex: sync.Mutex{}, Dvr: msgProtobuf.GetDistanceVectorRouting()}
						presence.HandleRouting(session, cnf, stream, nl, rt, otherDvr, timeTook)
						break
					case protobuf.RequestType_RETRANSMIT:
						msgProtobuf.Target = msgProtobuf.Sender
						msgProtobuf.Sender = cnf.NodeName

						msgProtobuf.Content = &protobuf.Header_ServerVideoChunk{
							ServerVideoChunk: &protobuf.ServerVideoChunk{
								Data: []byte(msgProtobuf.GetClientCommand().AdditionalInformation + "Hi there"),
							},
						}

						// send it back
						data, err := proto.Marshal(&msgProtobuf)
						if err != nil {
							log.Println("Error marshal protobuf:", err)
							break
						}
						err = config.SendMessage(stream, data)
						if err != nil {
							log.Println("Error sending message:", err)
							break
						}
					}

				}(stream)

			}
		}()
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
