package pubsub

import "github.com/DemonDCC/pubsub/router"

// Broker -
type Broker interface {
	Topic

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

	Subscribe(topic string, handler router.Handler) (Subscriber, error)
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
