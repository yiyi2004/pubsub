package pubsub

// Subscription -
type Subscription interface {
	Subscribe(topic string, handler Handler) error
	SubscribeChan(topic string, channel chan Msg) error

	Unsubscribe() (int, error)
	AutoUnsubcribeAll(max int) error

	Err() error

	Close() error
}

// Handler -
type Handler interface{}
