package client

import (
	"fmt"
	"log"
	"net"
	"time"

	bootstrapsystem "github.com/rodrigo0345/esr-tp2/client/bootstrapSystem"
	keepalive "github.com/rodrigo0345/esr-tp2/client/keep-alive"
	"github.com/rodrigo0345/esr-tp2/client/ui"
	"github.com/rodrigo0345/esr-tp2/config"
	cnf "github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
)

func Client(cnf *cnf.AppConfigList) {
	var bsAddr string
	if cnf.Neighbors == nil {
		bsAddr = "10.0.13.10:3242"
	} else {
		bsAddr = fmt.Sprintf("%s:%d", cnf.Neighbors[0].Ip, cnf.Neighbors[0].Port)
	}

	logger := config.NewLogger(2, cnf.NodeName)
	bsSystem := bootstrapsystem.NewBsSystem(logger, cnf)
	go bsSystem.ProcessBsCommands()

	laddr := net.UDPAddr{
		Port: 2222,
		IP:   net.IPv4zero, 
	}
	listener, err := net.ListenUDP("udp", &laddr)
	if err != nil {
		log.Fatalf("Failed to create UDP listener: %v", err)
	}
	defer listener.Close()

	createUIChan := make(chan *bootstrapsystem.Pop)

	uiChannel := make(chan *protobuf.Header, 100)

	go RequestNeighbors(cnf, bsAddr, bsSystem)

	go func() {
		for {
			data, _, err := config.ReceiveMessageUDP(listener)
			if err != nil {
				log.Printf("Error receiving message: %s\n", err)
				continue
			}

			go handleMessage(data, cnf, bsSystem, logger, createUIChan, uiChannel)
		}
	}()

	// Wait until we receive the bestPop to start the UI
	select {
	case bestPop := <-createUIChan:
    logger.Info(fmt.Sprintf("Best pop is %s", bestPop.Addr))
    go keepalive.KeepAlive(cnf, bestPop.Addr, 4)
		ui.StartUI(bestPop.Addr, cnf, uiChannel)
	}
}

func handleMessage(data []byte, cnf *cnf.AppConfigList, bsSystem *bootstrapsystem.BootstrapSystem, logger *config.Logger, createUIChan chan<- *bootstrapsystem.Pop, uiChannel chan<- *protobuf.Header) {
	msg, err := config.UnmarshalHeader(data)
	if err != nil {
		log.Printf("Error unmarshaling message: %s\n", err)
		return
	}

	switch msg.Type {
	case protobuf.RequestType_CLIENT_PING:
		if msg.Sender == "bs" {
			logger.Debug("Received Neighbors from bootstrapper")
			cbNb := ProcessNeighbors(cnf, msg, bsSystem)

			select {
			case result := <-cbNb:
				if result.Success {
					logger.Debug("Successfully received neighbors from bootstrapper")
					ProcessBestPop(cnf, msg, bsSystem)
					return
				}
				logger.Debug("Failed to receive neighbors from bootstrapper")
			}
			return
		}

		logger.Debug("Received POP Ping")
		callback := make(chan *bootstrapsystem.CallbackData)
		bsSystem.Command <- &bootstrapsystem.BsCommand{
			Command: bootstrapsystem.PROCESS_PING,
			Data: &bootstrapsystem.ChannelData{
				Msg:      msg,
				Callback: callback,
			},
		}

		select {
		case data := <-callback:
			if data.Success {
				bestPop := bsSystem.GetBestPop()
				// Signal to create the UI if not already signaled
				select {
				case createUIChan <- bestPop:
				default:
					// Already signaled
				}
			} else {
				panic("Error bootstrapping")
			}
		}
	case protobuf.RequestType_RETRANSMIT:
		logger.Debug("Received video chunk")
		// Send the message to the UI channel
		select {
		case uiChannel <- msg:
		default:
			logger.Error("UI channel is not ready to receive messages")
		}
	default:
		// Handle other message types if necessary
	}
}

func RequestNeighbors(cnf *cnf.AppConfigList, bsAddr string, bsSystem *bootstrapsystem.BootstrapSystem) {
	for {
		callback := make(chan *bootstrapsystem.CallbackData)
		bsSystem.Command <- &bootstrapsystem.BsCommand{
			Command: bootstrapsystem.REQUEST_NEIGHBORS,
			Data: &bootstrapsystem.ChannelData{
				ContactAddr: bsAddr,
				Callback:    callback,
			},
		}

		select {
		case <-callback:
		}

		if bsSystem.Neighbors != nil {
			break
		}
		time.Sleep(time.Second * 3)
	}
}

func ProcessNeighbors(cnf *cnf.AppConfigList, msg *protobuf.Header, bsSystem *bootstrapsystem.BootstrapSystem) chan *bootstrapsystem.CallbackData {
	callback := make(chan *bootstrapsystem.CallbackData)
	bsSystem.Command <- &bootstrapsystem.BsCommand{
		Command: bootstrapsystem.PROCESS_NEIGHBORS,
		Data: &bootstrapsystem.ChannelData{
			Msg:      msg,
			Callback: callback,
		},
	}
	return callback
}

func ProcessBestPop(cnf *cnf.AppConfigList, msg *protobuf.Header, bsSystem *bootstrapsystem.BootstrapSystem) {
	for {
		callback := make(chan *bootstrapsystem.CallbackData)
		bsSystem.Command <- &bootstrapsystem.BsCommand{
			Command: bootstrapsystem.REQUEST_PING,
			Data: &bootstrapsystem.ChannelData{
				Msg:      msg,
				Callback: callback,
			},
		}
		select {
		case data := <-callback:
			if data.Success {
			}
		}

		if bsSystem.BestPop != nil {
			break
		}
		time.Sleep(time.Second * 3)
	}
}
