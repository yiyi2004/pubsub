package pubsub

// Subscriber -
type Subscriber interface {
	SubscribeChan(b Broker, topic string, ch interface{}) error

	Unsubscribe() (int, error)
	AutoUnsubcribeAll(max int) error

	Err() error

	Close() error
}

// Handler -
type Handler interface{}
