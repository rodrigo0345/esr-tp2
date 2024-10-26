package presence

import (
	"fmt"
	"log"
	"github.com/rodrigo0345/esr-tp2/config"
	sendmessage "github.com/rodrigo0345/esr-tp2/config/send_message"
	"google.golang.org/protobuf/proto"
)

func PingNeighbors(cnf *config.AppConfigList) NeighborList {
	fmt.Printf("Sincronizing neighbors...\n")
	msg := sendmessage.Message{
		Target:  "",
		Content: *cnf.NodeIp,
		Type:    sendmessage.TypeInteraction_NEIGHBOR,
	}
	var nl NeighborList
	for _, neighbor := range cnf.Neighbors {
		fmt.Printf("Pinging %s\n", neighbor)

		msg.Target = neighbor

		data, err := proto.Marshal(&msg)
		if err != nil {
			log.Fatal("Error marshling ping")
		}

		_, err = config.SendMessage(neighbor, data)
		if err != nil {
			log.Println("Error sending ping")
			continue
		}

		nl.AddNeighbor(&Interface{
			Ip:   neighbor,
			Port: 4242,
		})
	}
	return nl
}

func (n *NeighborList) AddNeighbor(neighbor *Interface) {
	n.content = append(n.content, *neighbor)
}
