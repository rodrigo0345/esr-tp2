package server

import (
	"github.com/rodrigo0345/esr-tp2/config"
)

func Server(config *config.AppConfigList) {

  // get rid of / and .mp4
  client := Client{}
  (&client).NewClient(*config.VideoUrl, config.ServerUrl.Port, config.ServerUrl.Url)
  (&client).SendSource()
}
