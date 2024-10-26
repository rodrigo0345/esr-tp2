package config

import "github.com/rodrigo0345/esr-tp2/config/protobuf"

// enum with Server, Support, Client

type AppConfig uint8

const (
	Server AppConfig = iota
	Presence
	Client
)

type ServerUrl struct {
  Url string
  Port int
}

type AppConfigList struct {
	Topology     AppConfig
	VideoUrl     *string
	Neighbors []*protobuf.Interface
  NodeIP    *protobuf.Interface
  NodeIp *protobuf.Interface
}
