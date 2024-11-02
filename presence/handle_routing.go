package presence

import (
	"fmt"
	"log"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	dvrouting "github.com/rodrigo0345/esr-tp2/presence/distance_vector_routing"
	"google.golang.org/protobuf/proto"
)

func HandleRouting(conn quic.Connection, cnf *config.AppConfigList, stream quic.Stream, nl *NeighborList, dvr *dvrouting.DistanceVectorRouting, otherDvr *dvrouting.DistanceVectorRouting, timeTook int32) {

	// the source needs to be updated only when sending the routing table
	remote := conn.RemoteAddr().String()

	// send our routing table back
	msg := protobuf.Header{
		Type:      protobuf.RequestType_ROUTINGTABLE,
		Length:    0,
		Timestamp: int32(time.Now().UnixMilli()),

		Content: &protobuf.Header_DistanceVectorRouting{
			DistanceVectorRouting: dvr.Dvr,
		},
	}
	msg.Length = int32(proto.Size(&msg))

	local := conn.LocalAddr().String()
	in := config.ToInterface(local)
	in.Port = cnf.NodeIP.Port

	msg.Sender = fmt.Sprintf("%v:%v", in.Ip, in.Port)

	rm := config.ToInterface(remote)
	remoteIp := rm.Ip
  remotePort := otherDvr.Dvr.Source.Port

	otherDvr.UpdateSource(dvrouting.Interface{Interface: &protobuf.Interface{
		Ip:   remoteIp,
		Port: remotePort,
	}})

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

	in.Port = otherDvr.Dvr.Source.Port

	dvr.WeakUpdate(cnf, otherDvr, timeTook)
	nl.AddNeighbor(in, cnf)
}
