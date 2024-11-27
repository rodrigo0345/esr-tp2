package bootstrapsystem

import (
	"fmt"
	"time"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
)

type Command string

const (
	REQUEST_NEIGHBORS Command = "REQUEST_NEIGHBORS"
	PROCESS_NEIGHBORS Command = "PROCESS_NEIGHBORS"
	REQUEST_PING      Command = "REQUEST_PING"
	PROCESS_PING      Command = "PROCESS_PING"
)

type BsCommand struct {
	Command Command
	Data    *ChannelData
}

type BootstrapSystem struct {
	Command chan *BsCommand
	Logger   *config.Logger
	Config   *config.AppConfigList
	BestPop  *Pop
  Neighbors []string
}

type ChannelData struct {
  Msg *protobuf.Header
	Callback chan *CallbackData
  ContactAddr string
}

type CallbackData struct {
	Success bool
	Error   error
}

type Pop struct {
	Addr       string
	AvgLatency time.Duration
}

func NewBsSystem(logger *config.Logger, config *config.AppConfigList) *BootstrapSystem {

	return &BootstrapSystem{
		Command: make(chan *BsCommand),
		Logger:   logger,
		Config:   config,
    BestPop: nil,
	}
}

func (bs *BootstrapSystem) ProcessBsCommands() {
	for {
		select {
		case command := <-bs.Command:
      bs.Logger.Debug(fmt.Sprintf("Received command: %v", command.Command))
			switch command.Command {
			case REQUEST_NEIGHBORS:
				go bs.ReqNeighbors(command.Data)
			case PROCESS_NEIGHBORS:
				go bs.ProcessNeighbors(command.Data)
			case REQUEST_PING:
				go bs.ReqPing(command.Data)
			case PROCESS_PING:
				go bs.ProcessPing(command.Data)
			}
		}
	}
}

func (bs *BootstrapSystem) GetBestPop() *Pop {
  if bs.BestPop == nil {
    panic("BestPop is not set")
  }
  return bs.BestPop
}
