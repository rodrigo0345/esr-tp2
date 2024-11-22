package presence

import (
	"sync"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
)

type NeighborList struct {
	mutex   sync.Mutex
	Content []*protobuf.Interface
}

func (nl *NeighborList) Has(neighbor *protobuf.Interface) bool {
	nl.mutex.Lock()
	defer nl.mutex.Unlock()

	for _, n := range nl.Content {
		if n.String() == neighbor.String() {
			return true
		}
	}

	return false
}

func (nl *NeighborList) Add(neighbor *protobuf.Interface, cnf *config.AppConfigList) {
  nl.mutex.Lock()
  defer nl.mutex.Unlock()

	if nl.HasNeighbor(neighbor) {
		return
	}

	if neighbor.Ip == cnf.NodeIP.Ip && neighbor.Port == cnf.NodeIP.Port {
		return
	}
	nl.Content = append(nl.Content, neighbor)
}
