package presence

import (
	"fmt"
	"log"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	dvr "github.com/rodrigo0345/esr-tp2/presence/distance_vector_routing"

	"google.golang.org/protobuf/proto"
)

func HandleRetransmit(session quic.Connection, cnf *config.AppConfigList, stream quic.Stream, neighborList *NeighborList, routingTable *dvr.DistanceVectorRouting, data *protobuf.Header, neighborsConnectionsMap *dvr.ConnectionPool) {
	target := data.GetTarget()
	nextHop, err := routingTable.GetNextHop(target)

	if err != nil {
		log.Printf("No next hop found: %v\n", err)
		return
	}

	// open a connection with nextHop and try to send the message
	// if it fails, send to another neighbor, and remove it from the list
	neighbor := nextHop.NextNode
	nextNeighbor := 0

	// propagate the message to the next neighbor
	for {
		var msg []byte
		neighborStream, _, err := neighborsConnectionsMap.GetConnectionStream(neighbor.String())
		defer config.CloseStream(neighborStream)

		if err != nil {
			log.Printf("Error starting stream to %s: %v\n", nextHop.String(), err)
			goto fail
		}

		msg, err = proto.Marshal(data)
		err = config.SendMessage(neighborStream, msg)
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
			goto fail
		}

		// wait to get a confirmation from the other side
		_, err = config.ReceiveMessage(neighborStream)
		if err != nil {
			log.Printf("Error receiving message: %v\n", err)
			goto fail
		}

		break

	fail:
		// try with another neighbor
		neighbor = neighborList.Content[nextNeighbor]
		nextNeighbor++
		continue
	}

}

func HandleRetransmitFromClient(cnf *config.AppConfigList, neighborList *NeighborList, routingTable *dvr.DistanceVectorRouting, data *protobuf.Header, neighborsConnectionsMap *dvr.ConnectionPool) {
	target := data.GetTarget()
	nextHop, err := routingTable.GetNextHop(target)

	if err != nil {
		log.Printf("Error getting next hop: %v\n", err)
		return
	}

	// open a connection with nextHop and try to send the message
	// if it fails, send to another neighbor, and remove it from the list
	neighbor := nextHop.NextNode
	fmt.Printf("Next hop: %v\n", nextHop.NextNode)
	nextNeighbor := 0

	// propagate the message to the next neighbor
	for {
		var msg []byte

		nIPString := fmt.Sprintf("%s:%d", neighbor.Ip, neighbor.Port)
		fmt.Printf("Sending message to %s", nIPString)
		neighborStream, _, err := neighborsConnectionsMap.GetConnectionStream(nIPString)
		defer config.CloseStream(neighborStream)

		if err != nil {
			log.Printf("Error starting stream to %s: %v\n", nextHop.String(), err)
			goto fail
		}

		msg, err = proto.Marshal(data)
		err = config.SendMessage(neighborStream, msg)
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
			goto fail
		}

		// wait to get a confirmation from the other side
		_, err = config.ReceiveMessage(neighborStream)
		if err != nil {
			log.Printf("Error receiving message: %v\n", err)
			goto fail
		}

		break

	fail:
		// try with another neighbor
		neighbor = neighborList.Content[nextNeighbor]
		nextNeighbor++
		continue
	}
}
