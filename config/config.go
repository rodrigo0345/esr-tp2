package config

// enum with Server, Support, Client

type AppConfig uint8

const (
	Server AppConfig = iota
	Support
	Client
)

type ServerUrl struct {
  Url string
  Port int
}

type AppConfigList struct {
	Topology     AppConfig
	VideoUrl     *string
	BootstrapUrl *string
  ServerUrl    *ServerUrl
}
