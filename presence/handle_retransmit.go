package presence

import (
	"log"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	distancevectorrouting "github.com/rodrigo0345/esr-tp2/presence/distance_vector_routing"
	"google.golang.org/protobuf/proto"
)

func HandleRetransmit(session quic.Connection, cnf *config.AppConfigList, stream quic.Stream, neighborList *NeighborList, routingTable *distancevectorrouting.DistanceVectorRouting, data *protobuf.Header) {
	target := data.GetTarget()
	nextHop, err := routingTable.GetNextHop(target)

	if err != nil {
		log.Printf("Error getting next hop: %v\n", err)
		return
	}

	// open a connection with nextHop and try to send the message
	// if it fails, send to another neighbor, and remove it from the list
	connectionEstablished := false
	neighbor := nextHop.NextNode
	nextNeighbor := 0

	var newStream quic.Stream
	var conn quic.Connection
	for !connectionEstablished {
		newStream, conn, err = config.StartStream(neighbor.String())
		if err != nil {
			log.Printf("Error starting stream to %s: %v\n", nextHop.String(), err)
			// try with another neighbor
			neighbor = neighborList.content[nextNeighbor]
			nextNeighbor++
			continue
		}

		var msg []byte
		msg, err = proto.Marshal(data)
		err = config.SendMessage(newStream, msg)
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
			// try with another neighbor
			neighbor = neighborList.content[nextNeighbor]
			nextNeighbor++
			continue
		}

    // wait to get a confirmation from the other side
    _, err = config.ReceiveMessage(newStream)
    if err != nil {
			log.Printf("Error receiving message: %v\n", err)
			// try with another neighbor
			neighbor = neighborList.content[nextNeighbor]
			nextNeighbor++
			continue
		}

		connectionEstablished = true
	}

  defer conn.CloseWithError(0, "")
	defer config.CloseStream(newStream)

}
