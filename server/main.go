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
	if err := listenAndServe(tlsConfig, config.ServerUrl.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func listenAndServe(tls *tls.Config, serverPort int) error {
	port := fmt.Sprintf(":%d", serverPort)
	listener, err := quic.ListenAddr(port, tls, nil)
	if err != nil {
		return fmt.Errorf("failed to start QUIC listener: %w", err)
	}
	fmt.Printf("QUIC server is listening on port %d\n", serverPort)

	for {

		session, err := listener.Accept(context.Background())
		if err != nil {
			return fmt.Errorf("failed to accept session: %w", err)
		}

		go handleSession(session)
	}
}

func handleSession(session quic.Session) {
	for {
		stream, err := session.AcceptStream(context.Background())
		if err != nil {
			log.Println("failed to accept stream:", err)
			return
		}
		defer stream.Close()

		if err := processStream(stream); err != nil {
			log.Println("error processing stream:", err)
			return
		}
	}
}

func processStream(stream quic.Stream) error {
	for {
		lengthBuf := make([]byte, 4)
		if _, err := io.ReadFull(stream, lengthBuf); err != nil {
			return fmt.Errorf("failed to read length prefix: %w", err)
		}

		messageLength := binary.BigEndian.Uint32(lengthBuf)

		buf := make([]byte, messageLength)
		if _, err := io.ReadFull(stream, buf); err != nil {
			return fmt.Errorf("failed to read full message: %w", err)
		}

		chunk := &protobuf.VideoChunk{}
		if err := proto.Unmarshal(buf, chunk); err != nil {
			return fmt.Errorf("failed to unmarshal: %w", err)
		}

		log.Printf("Received video request for: %s", chunk.VideoId)

		if _, err := stream.Write([]byte("Hello, world!")); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}
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
