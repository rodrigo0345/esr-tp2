package streaming

import (
	"fmt"
	"net"
)

func RedirectPacketToClient(data []byte, clientIp string) {
	// open a connection with nextHop and try to send the message via UDP
	// Resolve the UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", clientIp)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	// Dial the UDP connection
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error dialing UDP connection:", err)
		return
	}
	defer conn.Close() // Close the connection when the function returns

	// Send the data
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	fmt.Println("Data sent successfully")

}
