package router

import pubsub "github.com/DemonDCC/pubsub/interface"

// HandlerFunc -
type HandlerFunc func(msg pubsub.Packet)

// HandlersChain -
type HandlersChain []HandlerFunc

// Handle -
func (h HandlerFunc) Handle(packet pubsub.Packet) {
	h(packet)
}

// Handler -s
type Handler interface {
	Handle(packet pubsub.Packet)
}

// Router -
type Router interface {
	IRoutes
	Group(string, ...HandlerFunc) *RouteGroup
}

// IRoutes -
type IRoutes interface {
	Use(middlewares ...HandlerFunc) IRoutes
}

// RouteGroup -
type RouteGroup struct {
}
