package main

import (
	"strconv"
	"wind-guide/config"
	"wind-guide/routes"

	"github.com/lesismal/arpc"
)

func main() {
	config.ReadConfig()

	server := arpc.NewServer()

	for index := range config.ConfigSettings.Routes {
		server.Handler.Handle(config.ConfigSettings.Routes[index].Path, routes.Handlers[config.ConfigSettings.Routes[index].Path])
	}

	server.Run(config.ConfigSettings.Server.Host + ":" + strconv.Itoa(config.ConfigSettings.Server.Port))
}
