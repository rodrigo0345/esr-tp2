package config

import (
	"log"
	"net"
	"strconv"

	"github.com/rodrigo0345/esr-tp2/config/protobuf"
)

// enum with Server, Support, Client

type AppConfig uint8

const (
	Server AppConfig = iota
	Presence
	Client
)

type ServerUrl struct {
	Url  string
	Port int
}

type AppConfigList struct {
	Topology  AppConfig
	VideoUrl  *string
	Neighbors []*protobuf.Interface
	NodeIP    *protobuf.Interface
	NodeName  string
}

// ToInterface converts a string address (ip:port) into a protobuf.Interface.
func ToInterface(addr string) *protobuf.Interface {
	ip, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatal(err)
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal(err)
	}
	return &protobuf.Interface{Ip: ip, Port: int32(p)}
}

