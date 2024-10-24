package presence

import (
	"fmt"
	"log"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/send_message"
	"google.golang.org/protobuf/proto"
)

type Interface struct {
	Ip   string
	Port int
}

// contains the best next node to reach the target
type RoutesTable struct {
	Routes map[string]Interface
}

func (rt *Interface) ToString() string {
	return fmt.Sprintf("%s:%d", rt.Ip, rt.Port)
}

type NeighborList struct {
	content []Interface
}

func Presence(config *config.AppConfigList) {
	// identify and alert neighbors
	neighborList := PingNeighbors(config)

	fmt.Printf("Node is running on %s\n", config.ServerUrl.Url)
	MainListen(config, &neighborList)

	// create a table of ips with the best interface for a given source
	// from time to time, update the table by trying to ping each source
	// in this message, send the
	// {
	//   "workingRoute": ["127.0.0.1"],
	//   "visited": [""]
	//   "target": "127.0.0.1:4005"
	// }
	//
	// once the message reaches the server, the server will send a message to the client
	// with the time it took to reach the target

}

func MainListen(cnf *config.AppConfigList, neighborList *NeighborList) {
	msg := config.ReceiveMessage(*cnf.NodeIp)
	// Unmarshal message
	var message send_message.Message
	err := proto.Unmarshal(msg, &message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Received message: %s\n", msg)

	switch message.Type {
	case send_message.TypeInteraction_FAILED:
		break
	case send_message.TypeInteraction_IS_ALIVE:
		break
	case send_message.TypeInteraction_NEIGHBOR:
		// TODO add the new neighbor to the neighbor list
    neighborList.AddNeighbor(&Interface{
      Ip: message.Content,
      Port: 4242,
    })
		break
	case send_message.TypeInteraction_TRANSMITION:
		break
	}

}
