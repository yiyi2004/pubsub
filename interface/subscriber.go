package pubsub

import pkg "github.com/DemonDCC/pubsub/packet"

// Subscriber -
type Subscriber interface {
	ChanSubscribe(b Broker, topic string, ch interface{}) error

	Unsubscribe() (int, error)
	AutoUnsubscribe(max int) error
}

// Then I want to us context.Context to manage goroutine

// HandleFunc -
type HandleFunc func(msg pkg.Packet)

// HandlerChain -
type HandlerChain []HandleFunc
