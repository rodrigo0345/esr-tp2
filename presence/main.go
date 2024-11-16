package presence

import (
	_ "encoding/json"
	"fmt"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"github.com/rodrigo0345/esr-tp2/presence/bootstrapper"
)

func Presence(cnf *config.AppConfigList) {

  // Boostrap the neighbors
  NeighborList, err := bootstrapper.BSGetNeighbors(cnf, cnf.Neighbors[0])

  if err != nil {
    fmt.Printf("Error getting neighbors: %s\n", err.Error())
    return
  }

  cnf.Neighbors = []*protobuf.Interface{}
  for _, neighbor := range NeighborList {
    fmt.Printf("Neighbor: %s\n", neighbor)
    nb := config.ToInterface(neighbor)
    cnf.Neighbors = append(cnf.Neighbors, nb)
  }

	presenceSystem := NewPresenceSystem(cnf)

  if err != nil {
    presenceSystem.Logger.Error(err.Error())
    return
  }

	go presenceSystem.HeartBeatNeighbors(7)

	// kill clients that dont ping in a while
	go presenceSystem.HeartBeatClients(5)

	presenceSystem.Logger.Info(fmt.Sprintf("Node is running on %s\n", presenceSystem.Config.NodeIP.String()))

	go presenceSystem.ListenForClients()
	presenceSystem.ListenForClientsInUDP()
}
