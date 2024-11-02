package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rodrigo0345/esr-tp2/client"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/config/protobuf"
	"github.com/rodrigo0345/esr-tp2/presence"
	"github.com/rodrigo0345/esr-tp2/server"
)

var (
	nodeName   = flag.String("n", "", "Name of the current node where clients can connect")
	specificIP = flag.String("s", "0.0.0.0:4242", "Specify an IP and port to bind to, default is 0.0.0.0:4242")
	neighbors  = flag.String("v", "", "Comma-separated list of neighbors in the format ip:port")
	debugMode  = flag.Bool("d", false, "Turn on debug mode")
)

func commandParser(command []string) (*config.AppConfigList, error) {
	if len(command) < 1 {
		return nil, errors.New("Invalid command, use main.go <server|presence|client> <arguments>")
	}

	flag.CommandLine.Parse(command[1:])

	switch command[0] {
	case "server":
		nodeIP := config.ToInterface(*specificIP)
		var neighborList []*protobuf.Interface
		if *neighbors != "" {
			neighborList = parseNeighbors(*neighbors)
		}
		config := config.AppConfigList{
			Topology:  config.Server,
			NodeIP:    nodeIP,
			Neighbors: neighborList,
		}
		return &config, nil
	case "presence":
		// Parse node IP and neighbors
		nodeIP := config.ToInterface(*specificIP)
		var neighborList []*protobuf.Interface
		if *neighbors != "" {
			neighborList = parseNeighbors(*neighbors)
		}
		config := config.AppConfigList{
			Topology:  config.Presence,
			NodeIP:    nodeIP,
			Neighbors: neighborList,
		}
		return &config, nil
	case "client":
		nodeIP := config.ToInterface(*specificIP)
		var neighborList []*protobuf.Interface
		if *neighbors != "" {
			neighborList = parseNeighbors(*neighbors)
		}
		config := config.AppConfigList{
			Topology:  config.Client,
			NodeIP:    nodeIP,
			Neighbors: neighborList,
		}
		return &config, nil
	default:
		return nil, errors.New("Invalid command, use main.go <server|presence|client> <arguments>")
	}
}

// parseNeighbors parses a comma-separated list of addresses into a slice of *protobuf.Interface
func parseNeighbors(neighborStr string) []*protobuf.Interface {
	neighborAddrs := strings.Split(neighborStr, ",")
	neighbors := make([]*protobuf.Interface, len(neighborAddrs))
	for i, addr := range neighborAddrs {
		neighbors[i] = config.ToInterface(strings.TrimSpace(addr))
	}
	return neighbors
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main.go <server|presence|client> <arguments>")
		os.Exit(1)
	}

	instConfig, err := commandParser(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	if *nodeName == "" {
		fmt.Println("Error: Node name not set, please use -n <node name>")
		os.Exit(1)
	}

	// Set the default IP and port if not provided
	if instConfig.NodeIP == nil {
		su := protobuf.Interface{Ip: "0.0.0.0", Port: 4242}
		instConfig.NodeIP = &su
	}
	instConfig.NodeName = *nodeName

	// Debug mode
	if *debugMode {
		log.Println("Debug mode is ON")
		log.Printf("Node Name: %s", *nodeName)
		log.Printf("Node IP: %s:%d", instConfig.NodeIP.Ip, instConfig.NodeIP.Port)
		for i, neighbor := range instConfig.Neighbors {
			log.Printf("Neighbor %d: %s:%d", i, neighbor.Ip, neighbor.Port)
		}
	}

	switch instConfig.Topology {
	case config.Server:
		server.Server(instConfig)
	case config.Presence:
		presence.Presence(instConfig)
	case config.Client:
		client.Client(instConfig)
	}
}
