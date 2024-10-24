package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"

	"github.com/quic-go/quic-go"
)

// SendMessage sends a message to the provided QUIC address and returns the sent message
func SendMessage(address string, message []byte) ([]byte, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // for testing only, don't use in production
		NextProtos:         []string{"quic-echo-example"},
	}

	session, err := quic.DialAddr(context.Background(), address, tlsConfig, nil)
	if err != nil {
		return nil, err
	}
	defer session.CloseWithError(0, "Client closed")

	// Open a stream
	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	_, err = stream.Write(message)
	if err != nil {
    return nil, err
	}

	fmt.Printf("Sent message: %x\n", message) // Printing message sent in hex
	return message, nil
}

func ReceiveMessage(address string) []byte {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // for testing only, don't use in production
		NextProtos:         []string{"quic-echo-example"},
	}

	// Dial the address using QUIC
	session, err := quic.DialAddr(context.Background(), address, tlsConfig, nil)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer session.CloseWithError(0, "Client closed")

	// Open a stream
	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		log.Fatalf("Failed to open stream: %v", err)
	}
	defer stream.Close()

	// Buffer to hold received data
	buf := make([]byte, 1024)
	var receivedMessage []byte

	for {
		// Read from the stream
		n, err := stream.Read(buf)
		if err != nil {
			if err == io.EOF {
				break // End of stream
			}
			log.Fatalf("Failed to read message: %v", err)
		}

		receivedMessage = append(receivedMessage, buf[:n]...)
		fmt.Printf("Received message (bytes): %x\n", buf[:n]) // Print received message in hex
	}

	return receivedMessage
}
