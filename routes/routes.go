package routes

import (
	"wind-guide/controllers"

	"github.com/lesismal/arpc"
)

var Handlers = map[string]func(ctx *arpc.Context){
	"/register-service":  controllers.RegisterService,
	"/discovery-service": controllers.DiscoveryService,
}
