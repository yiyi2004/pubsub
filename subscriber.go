package pubsub

// Subscriber -
type Subscriber interface {
	Subscribe()
}

// Subscription -
type Subscription interface {
	Subscribe(topic string, handler Handler) error
	SubscribeChan(topic string, channel chan Msg) error

	Unsubcribe(topic string) error
	UnsubscribeAll() (int, error)
	AutoUnsubscribe(topic string, max int) error
	AutoUnsubcribeAll(max int) error
}

// Handler -
type Handler interface{}
