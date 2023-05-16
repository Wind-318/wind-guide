// Package routes corresponds to the routes in the config file.
package routes

import (
	"wind-guide/controllers"

	"github.com/lesismal/arpc"
)

// Handlers is the map of routes and handlers.
var Handlers = map[string]func(ctx *arpc.Context){
	"/register-service":  controllers.RegisterService,
	"/discovery-service": controllers.DiscoveryService,
}
