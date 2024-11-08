package distancevectorrouting

import (
	"fmt"
	"math"
	"net"
	"sort"
	"strconv"
	"sync"

	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"google.golang.org/protobuf/proto"
)

type Interface struct {
	*protobuf.Interface
}

func (rt Interface) ToString() string {
	return net.JoinHostPort(rt.Ip, strconv.Itoa(int(rt.Port)))
}

// contains the best next node to reach the target
// this is a wrapper around the protobuf DistanceVectorRouting
type DistanceVectorRouting struct {
	Mutex sync.RWMutex
	Dvr   *protobuf.DistanceVectorRouting
}

func (dvr *DistanceVectorRouting) Lock() {
	dvr.Mutex.Lock()
}

func (dvr *DistanceVectorRouting) Unlock() {
	dvr.Mutex.Unlock()
}

// make a path to itself
func CreateDistanceVectorRouting(cnfg *config.AppConfigList) *DistanceVectorRouting {
	thisAddress := Interface{cnfg.NodeIP}
	dvr := &protobuf.DistanceVectorRouting{
		Entries: make(map[string]*protobuf.NextHop),
		Source:  thisAddress.Interface,
	}

	// add the path to itself with a distance of 0
	dvr.Entries[cnfg.NodeName] = &protobuf.NextHop{
		NextNode: thisAddress.Interface,
		Distance: 0,
	}
	return &DistanceVectorRouting{Mutex: sync.RWMutex{}, Dvr: dvr}
}

type NeighborRouting struct {
	Neighbor Interface
	Routing  DistanceVectorRouting
}

type RequestRoutingDelay struct {
	Neighbor Interface
	Delay    int64
}

func NewRouting(myName string, myIP *protobuf.Interface, neighborsRoutingTable []DistanceVectorRouting, requestRoutingDelay []RequestRoutingDelay) *DistanceVectorRouting {

	// Create a new DistanceVectorRouting instance to hold the combined routing table
	newRoutingTable := &DistanceVectorRouting{
		Mutex: sync.RWMutex{},
		Dvr: &protobuf.DistanceVectorRouting{
			Source:  myIP,
			Entries: make(map[string]*protobuf.NextHop),
		},
	}

	// Build a map of neighbor delays for quick lookup
	delayMap := make(map[string]int64)
	for _, req := range requestRoutingDelay {
		delayMap[req.Neighbor.ToString()] = req.Delay
	}

	// Iterate through each neighbor's distance vector table
	for _, neighborTable := range neighborsRoutingTable {
		neighborIP := Interface{neighborTable.Dvr.Source}.ToString()

		neighborDelay, _ := delayMap[neighborIP]

		for dest, nextHop := range neighborTable.Dvr.Entries {

			if dest == myName {
				newRoutingTable.Dvr.Entries[dest] = &protobuf.NextHop{
					NextNode: myIP, // Use the neighborID as the next node
					Distance: 0,
				}
			}

			var newDistance int64 = math.MaxInt64
			if nextHop.Distance != math.MaxInt64 {
				newDistance = nextHop.Distance + neighborDelay
			}

			existingNextHop, found := newRoutingTable.Dvr.Entries[dest]

			// Update if we found a shorter path or if the destination is new
			if !found || newDistance < existingNextHop.Distance {
				newRoutingTable.Dvr.Entries[dest] = &protobuf.NextHop{
					NextNode: neighborTable.Dvr.Source, // Use the neighborID as the next node
					Distance: newDistance,
				}
			}
		}
	}

	return newRoutingTable
}

func (dvr *DistanceVectorRouting) WeakUpdate(cnf *config.AppConfigList, other *DistanceVectorRouting, timeTook int64) {
	dvr.Lock()
	defer dvr.Unlock()

	for dest, nextHop := range other.Dvr.Entries {
		// Skip updating if the destination is this node
		if dest == cnf.NodeName {
			continue
		}

		// Calculate the new distance considering the time delay
		var newDistance int64 = math.MaxInt64
		if nextHop.Distance != math.MaxInt64 {
			newDistance = nextHop.Distance + timeTook
		}

		// Retrieve the current data in our routing table for this destination
		currentData, found := dvr.Dvr.Entries[dest]

		// Update if:
		// - The destination is not found (new route)
		// - The new route through the neighbor is shorter
		if !found || newDistance < currentData.Distance {
			dvr.Dvr.Entries[dest] = &protobuf.NextHop{
				NextNode: other.Dvr.Source,
				Distance: newDistance,
			}
		}
	}
}

func (dvr *DistanceVectorRouting) Remove(dest Interface) {
	delete(dvr.Dvr.Entries, dest.ToString())
}

func (dvr *DistanceVectorRouting) UpdateSource(source Interface) {
	dvr.Dvr.Source = source.Interface
}

func (dvr *DistanceVectorRouting) UpdateLocalSource(nodeName string, realPort int32, source Interface) {
	source.Port = realPort
	dvr.Dvr.Source = source.Interface
	dvr.Dvr.Entries[nodeName] = &protobuf.NextHop{
		NextNode: source.Interface,
		Distance: 0,
	}
}

func (dvr *DistanceVectorRouting) GetNextHop(dest string) (*protobuf.NextHop, error) {
	dvr.Mutex.RLock()
	defer dvr.Mutex.RUnlock()

	nextHop, found := dvr.Dvr.Entries[dest]
	if !found {
		return nil, fmt.Errorf("Destination %s not found", dest)
	}
	return nextHop, nil
}

func (dvr *DistanceVectorRouting) GetName(ip *protobuf.Interface) (string, error) {
	for key, nextHop := range dvr.Dvr.Entries {
		if fmt.Sprintf("%s:%d", nextHop.NextNode.Ip, nextHop.NextNode.Port) != fmt.Sprintf("%s:%d", ip.Ip, ip.Port) {
			continue
		}
		return key, nil
	}

	return "", nil
}

func (dvr *DistanceVectorRouting) Marshal() ([]byte, error) {
	return proto.Marshal(dvr.Dvr)
}

func (dvr *DistanceVectorRouting) Unmarshal(data []byte) error {
	return proto.Unmarshal(data, dvr.Dvr)
}

func (dvr *DistanceVectorRouting) Print() {
	fmt.Println("Routing Table:")

	// Coletar as entradas em uma slice para classificação
	entries := make([]string, 0, len(dvr.Dvr.Entries))
	for dest := range dvr.Dvr.Entries {
		entries = append(entries, dest)
	}

	// Ordenar as entradas por nome (destino)
	sort.Strings(entries)

	// Imprimir a tabela de roteamento ordenada
	for _, dest := range entries {
		nextHop := dvr.Dvr.Entries[dest]
		fmt.Printf("%s | %d | %s\n", dest, nextHop.Distance, Interface{nextHop.NextNode}.ToString())
	}
}
