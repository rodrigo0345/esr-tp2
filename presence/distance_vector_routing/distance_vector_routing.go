package distancevectorrouting

import (
	"fmt"
	"time"

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
    Source: thisAddress.Interface,
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

func NewRouting(myIP *protobuf.Interface, neighborsRoutingTable []DistanceVectorRouting, requestRoutingDelay []RequestRoutingDelay) *DistanceVectorRouting {
	// Create a new DistanceVectorRouting instance to hold the combined routing table
	newRoutingTable := &DistanceVectorRouting{
		&protobuf.DistanceVectorRouting{
			Source:  myIP,
			Entries: make(map[string]*protobuf.NextHop),
		},
	}

	// Build a map of neighbor delays for quick lookup
	delayMap := make(map[string]int)
	for _, req := range requestRoutingDelay {
		delayMap[req.Neighbor.ToString()] = req.Delay
	}

	// Iterate through each neighbor's distance vector table
	for _, neighborTable := range neighborsRoutingTable {
		neighborIP := Interface{neighborTable.Source}.ToString()

		neighborDelay, exists := delayMap[neighborIP]
		if !exists {
			// If no delay info exists for this neighbor, continue
			continue
		}

		// Iterate through the routing entries in the neighbor's table
		for dest, nextHop := range neighborTable.Entries {
			// Calculate the potential new distance to the destination through this neighbor
			newDistance := nextHop.Distance + int32(neighborDelay / 1000)

			existingNextHop, found := newRoutingTable.Entries[dest]

			// Update if we found a shorter path or if the destination is new
			if !found || newDistance < existingNextHop.Distance {
				newRoutingTable.Entries[dest] = &protobuf.NextHop{
					NextNode: neighborTable.Source, // Use the neighborID as the next node
					Distance: newDistance,
				}
			}
		}
	}

	return newRoutingTable
}

func (dvr *DistanceVectorRouting) WeakUpdate(other *DistanceVectorRouting, timeTook int32) {
	thisIP := Interface{dvr.Source}
	L := time.Since(time.UnixMilli(int64(timeTook)))

	for dest, nextHop := range other.Entries {
		// Skip updating if the destination is this node
		if dest == thisIP.ToString() {
			continue
		}

		// Calculate the new distance considering the time delay
		newDistance := nextHop.Distance + int32(L.Milliseconds())

		// Retrieve the current data in our routing table for this destination
		currentData, found := dvr.Entries[dest]

		// Update if:
		// - The destination is not found (new route)
		// - The new route through the neighbor is shorter
		if !found || newDistance < currentData.Distance {
			dvr.Entries[dest] = &protobuf.NextHop{
				NextNode: other.Source,
				Distance: newDistance,
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

func (dvr *DistanceVectorRouting) Print() {
  fmt.Println("Distance Vector Routing:")
  for dest, nextHop := range dvr.Entries {
    fmt.Printf("%s | %d | %s\n", dest, nextHop.Distance, Interface{nextHop.NextNode}.ToString())
  }
}
