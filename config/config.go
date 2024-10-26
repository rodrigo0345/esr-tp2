package config

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
	Neighbors []string
  ServerUrl    *ServerUrl
  NodeIp *string
}
