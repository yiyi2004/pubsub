package pubsub

// Subscriber -
type Subscriber interface {
	ChanSubscribe(b Broker, topic string, ch interface{}) error

	Unsubscribe() (int, error)
	AutoUnsubscribe(max int) error
}

// Then I want to us context.Context to manage goroutine

// Handler -
type Handler interface{}
