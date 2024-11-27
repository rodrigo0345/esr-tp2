package client

import (
	"fmt"
	"log"
	"net"
	"time"

	bootstrapsystem "github.com/rodrigo0345/esr-tp2/client/bootstrapSystem"
	"github.com/rodrigo0345/esr-tp2/client/ui"
	"github.com/rodrigo0345/esr-tp2/config"
	cnf "github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
)

var bsSystem *bootstrapsystem.BootstrapSystem
var updateUI func(*protobuf.Header) = nil

func Client(cnf *cnf.AppConfigList) {

	var bsAddr string
	if cnf.Neighbors == nil {
		bsAddr = "10.0.13.10:3242"
	} else {
		bsAddr = fmt.Sprintf("%s:%d", cnf.Neighbors[0].Ip, cnf.Neighbors[0].Port)
	}

	logger := config.NewLogger(2, cnf.NodeName)
	bsSystem = bootstrapsystem.NewBsSystem(logger, cnf)
	go bsSystem.ProcessBsCommands()

	listenIp := "0.0.0.0:2222"
	laddr := net.UDPAddr{
		Port: 2222,
		IP:   net.ParseIP(listenIp),
	}
	listener, err := net.ListenUDP("udp", &laddr)

	if err != nil {
		log.Fatalf("Failed to create UDP listener: %v", err)
	}

	go RequestNeighbors(cnf, bsAddr, bsSystem)

	// main listener loop
	func() {
		for {
			data, _, err := config.ReceiveMessageUDP(listener)

			if err != nil {
				log.Printf("Error receiving message: %s\n", err)
				continue
			}

			go func(data []byte) {

				msg, err := config.UnmarshalHeader(data)
				if err != nil {
					log.Printf("Error unmarshaling message: %s\n", err)
					return
				}

				switch msg.Type {
				case protobuf.RequestType_CLIENT_PING:

					// bootrapper logic
					// ping all the pops and find the best one
					if msg.Sender == "bs" {

						logger.Debug("Received Neighbors from bootstraper")
						cbNb := ProcessNeighbors(cnf, msg, bsSystem)

						select {
						case result := <-cbNb:

							if result.Success {
								logger.Debug("Successfully received neighbors from bootstraper")
								ProcessBestPop(cnf, msg, bsSystem)
							}

              logger.Debug("Failed to receive neighbors from bootstraper")
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
							updateUI = ui.StartUI(bestPop.Addr, cnf)
						} else {
							panic("Error bootstraping")
						}
						break
					}

					// find the best pop
					break
				case protobuf.RequestType_RETRANSMIT:

					if updateUI == nil {
						return
					}

					updateUI(msg)

				default:
				}
			}(data)
		}
	}()

}

func RequestNeighbors(cnf *cnf.AppConfigList, bsAddr string, bsSystem *bootstrapsystem.BootstrapSystem) {

  loop := true
	for loop {
		callback := make(chan *bootstrapsystem.CallbackData)
		bsSystem.Command <- &bootstrapsystem.BsCommand{
			Command: bootstrapsystem.REQUEST_NEIGHBORS,
			Data: &bootstrapsystem.ChannelData{
				ContactAddr: bsAddr,
				Callback:    make(chan *bootstrapsystem.CallbackData),
			},
		}

		select {
		case _ = <-callback:
		}

    if bsSystem.Neighbors != nil {
      loop = false
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
	callback := make(chan *bootstrapsystem.CallbackData)
	bsSystem.Command <- &bootstrapsystem.BsCommand{
		Command: bootstrapsystem.REQUEST_PING,
		Data: &bootstrapsystem.ChannelData{
			Msg:      msg,
			Callback: callback,
		},
	}
	select {
	case _ = <-callback:
	}
}
