package pubsub

import "context"

// Broker -
type Broker interface {
	Topics() []string
	NumTopics() int
	NumSubcribers(topic string) int
	Close() error

	AsyncSubscribe(ctx context.Context, topic string, handler HandlerFunc) (Subscription, error)
	SubscribeSync(ctx context.Context, topic string, handler HandlerFunc) (Subscription, error)
	QueueSubscribeSync(ctx context.Context, topic string, queue string) (Subscription, error)

	CreatePublisher(opts ...PublisherOptionFunc) Publisher
	CreateSubscription(opts ...SubscriptionOptionFunc) Subscription
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
type HandlerFunc func(in chan Packet, errChan chan error) (out chan Packet)

// HandlersChain -
type HandlersChain []HandlerFunc

// Last -
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}
