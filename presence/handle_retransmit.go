package presence

import (
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
)


func HandleRetransmit(ps *PresenceSystem, connection quic.Connection, stream quic.Stream, header *protobuf.Header) {
	target := header.GetTarget()

	nextHop, err := ps.RoutingTable.GetNextHop(target)
	if err != nil {
    ps.Logger.Error(err.Error())
		return
	}

	// open a connection with nextHop and try to send the message
	// if it fails, send to another neighbor, and remove it from the list
	neighbor := nextHop.NextNode
	nextNeighbor := 0

	// propagate the message to the next neighbor
	for {
		var msg []byte
		neighborIp := fmt.Sprintf("%s:%d", neighbor.Ip, neighbor.Port)
		neighborStream, _, err := ps.ConnectionPool.GetConnectionStream(neighborIp)
		defer config.CloseStream(neighborStream)

		if err != nil {
      ps.Logger.Error(err.Error())
			goto fail
		}

    msg, err = config.MarshalHeader(header)
		err = config.SendMessage(neighborStream, msg)

		if err != nil {
      ps.Logger.Error(err.Error())
			goto fail
		}

		// wait to get a confirmation from the other side
		_, err = config.ReceiveMessage(neighborStream)
		if err != nil {
      ps.Logger.Error(err.Error())
			goto fail
		}

		break

	fail:
		// try with another neighbor
		if len(ps.NeighborList.Content)-1 == nextNeighbor {
      ps.Logger.Error("No more neighbors to send the message")
			break
		}
		neighbor = ps.NeighborList.Content[nextNeighbor]
		nextNeighbor++
		continue
	}
}

func HandleRetransmitFromClient(ps *PresenceSystem, header *protobuf.Header) {
	target := header.GetTarget()
	nextHop, err := ps.RoutingTable.GetNextHop(target)

	if err != nil {
    ps.Logger.Error(err.Error())
		return
	}

	// open a connection with nextHop and try to send the message
	// if it fails, send to another neighbor, and remove it from the list
	neighbor := nextHop.NextNode
  ps.Logger.Info(fmt.Sprintf("Next hop: %v\n", nextHop.NextNode))
	nextNeighbor := 0

	// propagate the message to the next neighbor
	for {
		var msg []byte

		nIPString := fmt.Sprintf("%s:%d", neighbor.Ip, neighbor.Port)

    ps.Logger.Info(fmt.Sprintf("Sending message to %s\n", nIPString))

		neighborStream, _, err := ps.ConnectionPool.GetConnectionStream(nIPString)
		defer config.CloseStream(neighborStream)

		if err != nil {
      ps.Logger.Error(err.Error())
			goto fail
		}

		if neighborStream == nil {
      ps.Logger.Error(fmt.Sprintf("Error starting stream to %s: no connection found\n", nextHop.String()))
			goto fail
		}

    msg, err = config.MarshalHeader(header)
		err = config.SendMessage(neighborStream, msg)

		if err != nil {
      ps.Logger.Error(err.Error())
			goto fail
		}

		break

	fail:
		// try with another neighbor
		neighbor = ps.NeighborList.Content[nextNeighbor]
		nextNeighbor++
		continue
	}
}
