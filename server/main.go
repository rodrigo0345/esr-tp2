package server

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"github.com/rodrigo0345/esr-tp2/server/certs"
	"google.golang.org/protobuf/proto"
)

func Server(config *config.AppConfigList) {
	tls := generateTLS()
	listenAndServe(tls, config.ServerUrl.Port)
}

func listenAndServe(tls *tls.Config, serverPort int) error {
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
				defer stream.Close()

				for {
					// Read the 4-byte length prefix
					lengthBuf := make([]byte, 4)
					_, err := io.ReadFull(stream, lengthBuf)
					if err != nil {
						log.Println("Failed to read length prefix:", err)
						return
					}

					messageLength := binary.BigEndian.Uint32(lengthBuf)

					// Read the message of the specified length
					buf := make([]byte, messageLength)
					_, err = io.ReadFull(stream, buf)
					if err != nil {
						log.Println("Failed to read full message:", err)
						return
					}

					// Unmarshal the protobuf message
					chunk := &protobuf.VideoChunk{}
					err = proto.Unmarshal(buf, chunk)
					if err != nil {
						log.Println("Failed to unmarshal:", err)
						return
					}

					log.Println(chunk)

					// Send hello message to the client (optional)
					_, err = stream.Write([]byte("Hello, world!"))
					if err != nil {
						log.Println("Failed to write message:", err)
						return
					}
				}

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
