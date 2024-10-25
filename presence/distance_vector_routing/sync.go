// this file is used to sync all the nodes in the cluster with the routing table
package distancevectorrouting

import (
	"fmt"
	"sync"

	"github.com/rodrigo0345/esr-tp2/config/protobuf"
)

type Sync struct {
	sync.Mutex
	routingTable *DistanceVectorRouting
}

func (s *Sync) Update(routingTable *DistanceVectorRouting) {
	s.Lock()
	defer s.Unlock()
	s.routingTable = routingTable
}

func (s *Sync) GetRoutingTable() *DistanceVectorRouting {
	s.Lock()
	defer s.Unlock()
	return s.routingTable
}

func (s *Sync) Marshal() ([]byte, error) {
	s.Lock()
	defer s.Unlock()
	return s.routingTable.Marshal()
}

func (s *Sync) Unmarshal(data []byte) error {
	s.Lock()
	defer s.Unlock()
	return s.routingTable.Unmarshal(data)
}


