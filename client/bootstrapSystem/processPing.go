package bootstrapsystem

import (
	"fmt"
	"time"
)

func (bs *BootstrapSystem) ProcessPing(data *ChannelData) {
	header := data.Msg

	timestamp := time.Unix(header.Timestamp, 0)
	latency := time.Since(timestamp)

	addr := header.Sender

	if bs.BestPop == nil {
		bs.BestPop = &Pop{
			Addr:       addr,
			AvgLatency: latency,
		}
	} else {
		bs.Logger.Info(fmt.Sprintf("Best pop: %s, %d(ms)", addr, latency))
		if latency < bs.BestPop.AvgLatency {
			bs.BestPop.Addr = addr
			bs.BestPop.AvgLatency = latency
		}
	}

	data.Callback <- &CallbackData{
		Success: true,
		Error:   nil,
	}
}
