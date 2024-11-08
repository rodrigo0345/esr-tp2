package presence

import (
	_ "encoding/json"
	"fmt"
	"github.com/rodrigo0345/esr-tp2/config"
)

func Presence(config *config.AppConfigList) {

	presenceSystem := NewPresenceSystem(config)

	go presenceSystem.HeartBeatNeighbors(7)

	// kill clients that dont ping in a while
	go presenceSystem.HeartBeatClients(5)

	presenceSystem.Logger.Info(fmt.Sprintf("Node is running on %s\n", presenceSystem.Config.NodeIP.String()))

	go presenceSystem.ListenForClients()
	presenceSystem.ListenForClientsInUDP()
}
