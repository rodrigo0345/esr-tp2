package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf" // Import your generated protobuf package
	"google.golang.org/protobuf/proto"
)

func Client(config *config.AppConfigList) {
	// Address to send messages
	presenceNetworkEntrypointIp := config.Neighbors[0]
	pneIpString := fmt.Sprintf("%s:%d", presenceNetworkEntrypointIp.Ip, presenceNetworkEntrypointIp.Port)
	fmt.Println(pneIpString)

	listenIp := "127.0.0.1:2222"

	// Setup a UDP connection for sending
	conn, err := net.Dial("udp", pneIpString)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	// Create a UDP listener for incoming messages
	laddr := net.UDPAddr{
		Port: 2222, // Let the OS assign a port
		IP:   net.ParseIP(listenIp),
	}
	listener, err := net.ListenUDP("udp", &laddr)
	if err != nil {
		log.Fatalf("Failed to clone UDP listener: %v", err)
	}
	defer listener.Close()

	// Run a goroutine to listen for incoming messages
	go func() {
		for {
			buffer := make([]byte, 1024)
			n, _, err := listener.ReadFromUDP(buffer)
			if err != nil {
				log.Printf("Error receiving message: %s\n", err)
				continue
			}

      var msg protobuf.Header
      err = proto.Unmarshal(buffer[:n], &msg)
      if err != nil {
        log.Printf("Error unmarshaling message: %s\n", err)
        continue
      }

      // TODO: change all of this
			log.Printf("Received message from %s: %s\n", msg.Sender, msg.Target)
		}
	}()

	// Read messages from the terminal and send
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter messages to send. Type 'exit' to quit.")
	for {
		fmt.Print("Message: ")
		scanner.Scan()
		text := scanner.Text()
		if text == "exit" {
			break
		}
		fmt.Print("Target: ")
		scanner.Scan()
		target := scanner.Text()
		if text == "exit" {
			break
		}

		message := protobuf.Header{
			Type:      protobuf.RequestType_RETRANSMIT,
			Length:    0,
			Timestamp: int32(time.Now().UnixMilli()),
			ClientIp:  listenIp, // this is where the client is waiting for the response
			Sender:    "client",
			Target:    target,
			Content: &protobuf.Header_ClientCommand{
				ClientCommand: &protobuf.ClientCommand{
					Command:               protobuf.PlayerCommand_PLAY,
					AdditionalInformation: text,
				},
			},
		}

    fmt.Printf("Sending message: %s to %s\n", text, target)
		message.Length = int32(proto.Size(&message))
		data, err := proto.Marshal(&message)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			continue
		}

		_, err = conn.Write(data)
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		}
	}
}
