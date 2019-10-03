package pubsub

import "github.com/zhangce1999/pubsub/router"

// Broker -
type Broker interface {
	Topic
	Router

	CreatePublisher(opts ...PublisherOptionFunc) Publisher
	CreateSubscriber(opts ...SubscriberOptionFunc) Subscriber
}

// Topic -
type Topic interface {
	Topics() []string
	NumTopics() int
	RegisterTopic(topic string) (conn interface{}, err error)
	NumSubcribers(topic string) int
	Close(topics ...string)

	AsyncSubscribe(topic string, handler Handler) (Subscriber, error)
	Subscribe(topic string, handler Handler) (Subscriber, error)
}

// StatusInfo -
type StatusInfo interface {
	ConnStatus() ConnStatus

	// Channel Status

	// Server Status

	// Publisher Status

	// Subscription Status
}

// ConnStatus type
type ConnStatus int8

// Status -
const (
	DisConnected = ConnStatus(iota)
	Connected
	Closed
	Reconnecting
	Connecting

	DrainingSubs
	DrainingPubs
)

// HandlerFunc -
type HandlerFunc func(msg Packet)

// HandlersChain -
type HandlersChain []HandlerFunc

// Handle -
func (h HandlerFunc) Handle(packet Packet) {
	h(packet)
}

// Handler -s
type Handler interface {
	Handle(packet Packet)
}

// Router -
type Router interface {
	IRoutes
	Group(string, ...HandlerFunc) *router.Group
}

// IRoutes -
type IRoutes interface {
	Use(middlewares ...HandlerFunc) IRoutes
}
