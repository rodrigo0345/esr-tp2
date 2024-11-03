package config

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/quic-go/quic-go"
)

var timeout time.Duration = 400

func StartConnStream(address string) (quic.Stream, quic.Connection, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // for testing only, don't use in production
		NextProtos:         []string{"quic-echo-example"},
	}

	connection, err := quic.DialAddr(context.Background(), address, tlsConfig, &quic.Config{
		KeepAlivePeriod: timeout * time.Second,
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial: %w", err)
	}

	// Open a stream
	stream, err := connection.OpenStreamSync(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open stream: %w", err)
	}

	return stream, connection, nil
}

func CreateStream(conn quic.Connection) (quic.Stream, error) {
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %w", err)
	}
	return stream, nil
}

func CloseStream(stream quic.Stream) {
	if stream == nil {
		return
	}
	if err := stream.Close(); err != nil {
		log.Printf("Failed to close stream: %v", err)
	}
}

func CloseConnection(conn quic.Connection) {
	if err := conn.CloseWithError(0, "closing connection"); err != nil {
		fmt.Println("Failed to close connection:", err)
	}
}

func IsConnectionOpen(stream quic.Stream) bool {
	if stream == nil {
		return false
	}
	_, err := stream.Write([]byte("ping"))
	if err != nil {
		return false
	}
	return true
}

func SendMessage(stream quic.Stream, message []byte) error {
	// Create a context with a 10-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	if stream == nil {
		return fmt.Errorf("Sending data to a nil stream")
	}

	// Send the length of the message first
	lengthPrefix := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthPrefix, uint32(len(message)))

	// Wrap stream.Write in a channel to manage the timeout
	writeDone := make(chan error, 1)
	go func() {
		_, err := stream.Write(lengthPrefix)
		if err != nil {
			log.Printf("Error writing length prefix: %v\n", err)
			writeDone <- err
			return
		}
		_, err = stream.Write(message)
		writeDone <- err
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("send message timed out: %w", ctx.Err())
	case err := <-writeDone:
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}

	return nil
}

func ReceiveMessage(stream quic.Stream) ([]byte, error) {
	lengthBuf := make([]byte, 4)
	if err := readWithTimeout(stream, lengthBuf, timeout*time.Second); err != nil {
		return nil, fmt.Errorf("failed to read message length: %w", err)
	}

	messageLength := binary.BigEndian.Uint32(lengthBuf)
	if messageLength == 0 {
		return nil, fmt.Errorf("invalid message length: %d", messageLength)
	}

	buf := make([]byte, messageLength)
	if err := readWithTimeout(stream, buf, timeout*time.Second); err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	return buf, nil
}

// readWithTimeout reads data from the stream with a timeout.
func readWithTimeout(stream quic.Stream, buf []byte, timeout time.Duration) error {
	readDone := make(chan error, 1)

	go func() {
		_, err := io.ReadFull(stream, buf)
		readDone <- err
	}()

	select {
	case err := <-readDone:
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("read operation timed out after %v", timeout)
	}
}
