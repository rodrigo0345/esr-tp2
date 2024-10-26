package presence

import (
	"log"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	distancevectorrouting "github.com/rodrigo0345/esr-tp2/presence/distance_vector_routing"
	"google.golang.org/protobuf/proto"
)

func HandleRouting(stream quic.Stream, nl *NeighborList, dvr *distancevectorrouting.DistanceVectorRouting, otherDvr *distancevectorrouting.DistanceVectorRouting, timeTook int32) {

	// send our routing table back
	msg := protobuf.Header{
		Type:      protobuf.RequestType_ROUTINGTABLE,
		Length:    0,
		Timestamp: int32(time.Now().UnixMilli()),

		Content: &protobuf.Header_DistanceVectorRouting{
			DistanceVectorRouting: dvr.DistanceVectorRouting,
		},
	}
	msg.Length = int32(proto.Size(&msg))

	data, err := proto.Marshal(&msg)
	if err != nil {
		log.Printf("Error marshaling routing table: %v\n", err)
		return
	}

	err = config.SendMessage(stream, data)
	if err != nil {
		log.Printf("Error sending routing table: %v\n", err)
		return
	}

	// TODO: handle race condition
	dvr.WeakUpdate(otherDvr, timeTook)

	nl.AddNeighbor(otherDvr.Source)
}
