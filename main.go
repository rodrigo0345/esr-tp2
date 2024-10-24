package main

import (
	"errors"
	"log"
	"os"

	"github.com/rodrigo0345/esr-tp2/client"
	"github.com/rodrigo0345/esr-tp2/config"
	"github.com/rodrigo0345/esr-tp2/server"
	"github.com/rodrigo0345/esr-tp2/presence"
)


func commandParser(command []string) (*config.AppConfigList, error) {
	if command[0] == "-h" || command[0] == "--help" {
		return nil, errors.New(`
      Usage: main.go <server|support|client> <server->video_url>|<support->bootstrap_url>|<client->video_url>
      
      Commands:
        server: start server
        support: start support
        client: start client
      
      Arguments:
        --help, -h: show this help message

      Examples:
        main.go server ~/video.mp4
        main.go support http://localhost:8080
        main.go client http://localhost:8080/video.mp4
    `)
	}

	if len(command) < 2 {
		return nil, errors.New("Invalid command, use main.go <server|support|client> <server->video_url>|<support->bootstrap_url>|<client->video_url>")
	}
	switch command[0] {
	case "server":
		config := config.AppConfigList{Topology: config.Server, VideoUrl: &command[1]}
		return &config, nil
	case "presence":
    config := config.AppConfigList{Topology: config.Presence, ServerUrl: &command[1], Neighbors: command[2:]}
		return &config, nil
	case "client":
		config := config.AppConfigList{Topology: config.Client, VideoUrl: &command[1]}
		return &config, nil
	}

	return nil, errors.New("Invalid command, use main.go <server|support|client> <server->video_url>|<presence->nodeip->neighbors>|<client->video_url>")
}


func main() {
	instConfig, err := commandParser(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

  su := config.ServerUrl{Url: "127.0.0.1", Port: 4242}
  instConfig.ServerUrl = &su

	switch instConfig.Topology {
	case config.Server:
		server.Server(instConfig)
	case config.Support:
    presence.Presence(instConfig)
	case config.Client:
    client.Client(instConfig)
	}
}
