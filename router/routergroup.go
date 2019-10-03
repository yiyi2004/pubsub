package router

import pubsub "github.com/zhangce1999/pubsub/interface"

// Group is like gin.RouterGroup
type Group struct {
	root     bool
	basePath string
	engine   pubsub.Broker
	handlers pubsub.HandlersChain
}
