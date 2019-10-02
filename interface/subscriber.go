package pubsub

// Subscriber -
type Subscriber interface {
	Topic() string

	Unsubscribe() (int, error)
	AutoUnsubscribe(max int) error
}

// Then I want to us context.Context to manage goroutine

// // Subscription -
// type Subscription interface {
// 	Subscriber

// 	Topics() []string
// }
