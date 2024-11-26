package client

import (
	"fmt"
	"time"

  "github.com/go-ping/ping"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
)

func GetPops(cnf *config.AppConfigList, logger *config.Logger) ([]string, error) {

	serverIp := fmt.Sprintf("%s:%d", cnf.Neighbors[0].Ip, cnf.Neighbors[0].Port)

	msg := protobuf.Header{
		Type:           protobuf.RequestType_BOOTSTRAPER,
		Length:         0,
		Timestamp:      time.Now().UnixMilli(),
		ClientIp:       "nil",
		Sender:         "client",
		Target:         "s1",
		RequestedVideo: "nil",
		Content:        nil,
	}

	msg.Length = int32(proto.Size(&msg))
	msg.Target = "s1"

	data, err := proto.Marshal(&msg)

	if err != nil {
		return nil, err
	}

	stream, _, err := config.StartConnStream(serverIp)
	defer config.CloseStream(stream)

	if err != nil {
		return nil, err
	}

	config.SendMessage(stream, data)

  data, err = config.ReceiveMessage(stream)
  if err != nil {
    return nil, err
  }

  proto.Unmarshal(data, &msg)
  neighbors := msg.GetBootstraperResult().GetNeighbors()

  logger.Info(fmt.Sprintf("Neighbors: %v\n", neighbors))

	return neighbors, nil
}


func GetBestPop(pops []string, logger *config.Logger) BestPop {
	var bestPop BestPop
	bestPop.Avg = time.Hour // Initialize to a very high value

	for _, pop := range pops {
		popIp := config.ToInterface(pop)
		pinger, err := ping.NewPinger(popIp.Ip)

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to create pinger for %s: %v", pop, err))
			continue
		}

		pinger.Count = 1
		pinger.Timeout = 5 * time.Second
		pinger.SetPrivileged(true)

		err = pinger.Run()
		if err != nil {
			logger.Error(fmt.Sprintf("Ping failed for %s: %v", pop, err))
			continue
		}

		stats := pinger.Statistics() // Get the statistics
		if stats.AvgRtt < bestPop.Avg {
			bestPop.Avg = stats.AvgRtt
			bestPop.Ip = popIp
		}
	}

	logger.Info(fmt.Sprintf("Best pop: %v, Avg ping: %v", bestPop.Ip, bestPop.Avg))

	return bestPop
}
