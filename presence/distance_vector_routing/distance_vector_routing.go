package distancevectorrouting

import (
	"fmt"

	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
)

type Interface struct {
  *protobuf.Interface
}

func (rt Interface) ToString() string {
	return fmt.Sprintf("%s:%d", rt.Ip, rt.Port)
}

// contains the best next node to reach the target
// this is a wrapper around the protobuf DistanceVectorRouting
type DistanceVectorRouting struct {
	*protobuf.DistanceVectorRouting
}

// make a path to itself
func CreateDistanceVectorRouting(thisAddress Interface) *DistanceVectorRouting {
	dvr := &protobuf.DistanceVectorRouting{
		Entries: make(map[string]*protobuf.NextHop),
  }

	// add the path to itself with a distance of 0
	dvr.Entries[thisAddress.ToString()] = &protobuf.NextHop{
		NextNode: thisAddress.Interface,
		Distance: 0,
	}
	return &DistanceVectorRouting{dvr}
}

type NeighborRouting struct {
	Neighbor Interface
	Routing  DistanceVectorRouting
}

type RequestRoutingDelay struct {
	Neighbor Interface
	Delay    int
}

func (dvr *DistanceVectorRouting) Update(neighborsRouting []NeighborRouting, requestRoutingDelay []RequestRoutingDelay) {
	// Build a map of neighbor delays for quick lookup
	delayMap := make(map[Interface]int)
	for _, req := range requestRoutingDelay {
		delayMap[req.Neighbor] = req.Delay
	}

	for _, neighbor := range neighborsRouting {

		neighborDelay, exists := delayMap[neighbor.Neighbor]
		if !exists {
			// If no delay info exists for this neighbor, continue
			continue
		}

		for dest, nextHop := range neighbor.Routing.Entries {

			// Calculate the potential new distance to the destination through this neighbor
			newDistance := nextHop.Distance + int32(neighborDelay)
			existingNextHop, found := dvr.Entries[dest]

			// Update if we found a shorter path or if the destination is new
			if !found || newDistance < existingNextHop.Distance {
				dvr.Entries[dest] = &protobuf.NextHop{
					NextNode: neighbor.Neighbor.Interface,
					Distance: newDistance,
				}
			}
		}
	}
}

func (dvr *DistanceVectorRouting) Remove(dest Interface) {
	delete(dvr.Entries, dest.ToString())
}

func (dvr *DistanceVectorRouting) GetNextHop(dest Interface) (*protobuf.NextHop, error) {
	nextHop, found := dvr.Entries[dest.ToString()]
	if !found {
		return nil, fmt.Errorf("Destination %s not found", Interface(dest).ToString())
	}
	return nextHop, nil
}

func (dvr *DistanceVectorRouting) Marshal() ([]byte, error) {
  return proto.Marshal(dvr)
}

func (dvr *DistanceVectorRouting) Unmarshal(data []byte) error {
  return proto.Unmarshal(data, dvr)
}
