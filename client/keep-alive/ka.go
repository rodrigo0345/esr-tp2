package keepalive

import (
	"time"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
)

func KeepAlive(cnf *config.AppConfigList, bestPopAddr string, timeout time.Duration) {
	for {
		header := &protobuf.Header{
			Type:       protobuf.RequestType_HEARTBEAT,
			Length:     0,
			Timestamp:  int64(time.Now().Unix()),
			Sender:     cnf.NodeIP.Ip,
			Target:     []string{bestPopAddr},
			ClientPort: "2222", // the port where this listens
		}

		data, err := proto.Marshal(header)

		if err != nil {
			continue
		}

		config.SendMessageUDP(bestPopAddr, data)
		time.Sleep(timeout * time.Second)
	}
}
