package pubsub

// Subscriber -
type Subscriber interface {
	Topic() string
	Drain()
	Delivered() (int64, error)

	Unsubscribe() (int, error)
	AutoUnsubscribe(max int) error

	MaxPendings(int, int, error)
	PendingLimits(int, int, error)
	SetPendingLimits(msgLimit, bytesLimit int) error
}

// Then I want to us context.Context to manage goroutine

// // Subscription -
// type Subscription interface {
// 	Subscriber

// 	Topics() []string
// }
