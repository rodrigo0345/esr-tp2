package server

import (
	"github.com/rodrigo0345/esr-tp2/config"
)

func Server(cnf *config.AppConfigList) {
	serverSystem := NewServerSystem(cnf)
	go serverSystem.HeartBeatNeighbors(5)
	go serverSystem.BackgroundStreaming()
  go serverSystem.BsClients()
	serverSystem.ListenForClients()
}
