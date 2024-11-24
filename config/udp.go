package config

import (
	"fmt"
	"net"
)

// SendMessageUDP sends a message to the given UDP address with length-prefix encoding
func SendMessageUDP(address string, message []byte) error {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return fmt.Errorf("failed to dial UDP address %s: %w", address, err)
	}
	defer conn.Close()

	// Send the message directly without a length prefix.
	_, err = conn.Write(message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// ReceiveMessageUDP listens on the given UDP connection and receives a single message
func ReceiveMessageUDP(conn *net.UDPConn) ([]byte, *net.UDPAddr, error) {
	// Allocate a buffer to hold the incoming datagram.
	// 65535 bytes is the maximum size for a UDP packet.
	buf := make([]byte, 65535)

	// Read a single datagram from the connection.
	n, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read message: %w", err)
	}

	if n == 0 {
		return nil, nil, fmt.Errorf("received empty message")
	}

	// Return the slice containing the actual message data.
	return buf[:n], addr, nil
}
