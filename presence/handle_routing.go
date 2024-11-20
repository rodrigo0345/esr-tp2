package presence

import (
	"fmt"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	dvrouting "github.com/rodrigo0345/esr-tp2/presence/distance_vector_routing"
	"google.golang.org/protobuf/proto"
)

func HandleRouting(ps *PresenceSystem, conn quic.Connection, stream quic.Stream, header *protobuf.Header) {

	// the source needs to be updated only when sending the routing table
	remote := conn.RemoteAddr().String()
	receivedDvr := &dvrouting.DistanceVectorRouting{Mutex: sync.RWMutex{}, Dvr: header.GetDistanceVectorRouting()}
	timeTook := time.Since(time.UnixMilli(int64(header.Timestamp)))

	// send our routing table back
	msg := protobuf.Header{
		Type:      protobuf.RequestType_ROUTINGTABLE,
		Length:    0,
		Timestamp: time.Now().UnixMilli(),

		Content: &protobuf.Header_DistanceVectorRouting{
			DistanceVectorRouting: ps.RoutingTable.Dvr,
		},
	}

	msg.Length = int32(proto.Size(&msg))

	local := conn.LocalAddr().String()
	in := config.ToInterface(local)
	in.Port = ps.Config.NodeIP.Port

	msg.Sender = fmt.Sprintf("%v:%v", in.Ip, in.Port)

	rm := config.ToInterface(remote)
	remoteIp := rm.Ip
	remotePort := receivedDvr.Dvr.Source.Port

	// update the source of the received dvr
	receivedDvr.UpdateSource(dvrouting.Interface{Interface: &protobuf.Interface{
		Ip:   remoteIp,
		Port: remotePort,
	}})

	data, err := proto.Marshal(&msg)
	if err != nil {
    ps.Logger.Error(err.Error())
		return
	}

	err = config.SendMessage(stream, data)
	if err != nil {
    ps.Logger.Error(err.Error())
		return
	}
  
  rm.Port = receivedDvr.Dvr.Source.Port

	ps.RoutingTable.WeakUpdate(ps.Config, receivedDvr, int64(timeTook))
	ps.NeighborList.AddNeighbor(rm, ps.Config)
}
