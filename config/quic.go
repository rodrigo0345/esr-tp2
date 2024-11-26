package config

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/quic-go/quic-go"
)

// StartConnStream initializes a QUIC connection and opens a stream
func StartConnStream(address string) (quic.Stream, quic.Connection, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}

	// Use a background context for dialing
	ctx := context.Background()

	connection, err := quic.DialAddr(ctx, address, tlsConfig, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial: %w", err)
	}

	// Open a stream without timeout
	stream, err := connection.OpenStreamSync(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open stream: %w", err)
	}

	return stream, connection, nil
}

// CreateStream opens a new stream in the provided connection
func CreateStream(conn quic.Connection) (quic.Stream, error) {
	ctx := context.Background()

	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %w", err)
	}
	return stream, nil
}

// CloseStream closes a given stream
func CloseStream(stream quic.Stream) {
	if stream == nil {
		return
	}
	if err := stream.Close(); err != nil {
		log.Printf("Failed to close stream: %v", err)
	}
}

// CloseConnection closes the given QUIC connection
func CloseConnection(conn quic.Connection) {
	if conn == nil {
		return
	}
	if err := conn.CloseWithError(0, "closing connection"); err != nil {
		log.Printf("Failed to close connection: %v", err)
	}
}

// IsConnectionOpen checks if the connection context has been canceled
func IsConnectionOpen(conn quic.Connection) bool {
	if conn == nil {
		return false
	}
	return conn.Context().Err() == nil
}

func IsStreamOpen(stream quic.Stream) bool {
  if stream == nil {
    return false
  }
  return stream.Context().Err() == nil
}

// SendMessage sends a message to the stream with length-prefix encoding
func SendMessage(stream quic.Stream, message []byte) error {
	if stream == nil {
		return fmt.Errorf("sending data to a nil stream")
	}

	// Send the length of the message first
	lengthPrefix := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthPrefix, uint32(len(message)))

	// Write length prefix
	if _, err := stream.Write(lengthPrefix); err != nil {
		return fmt.Errorf("failed to send message length: %w", err)
	}

	// Write the actual message
	if _, err := stream.Write(message); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// ReceiveMessage reads a message from the stream with length-prefix decoding
func ReceiveMessage(stream quic.Stream) ([]byte, error) {
	if stream == nil {
		return nil, fmt.Errorf("receive message from nil stream")
	}

	lengthBuf := make([]byte, 4)

	// Read the length of the message
	if _, err := io.ReadFull(stream, lengthBuf); err != nil {
		return nil, fmt.Errorf("failed to read message length: %w", err)
	}

	messageLength := binary.BigEndian.Uint32(lengthBuf)
	if messageLength == 0 {
		return nil, fmt.Errorf("invalid message length: %d", messageLength)
	}

	buf := make([]byte, messageLength)

	// Read the actual message
	if _, err := io.ReadFull(stream, buf); err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	return buf, nil
}

