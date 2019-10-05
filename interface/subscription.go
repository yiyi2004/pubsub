package pubsub

import "time"

// Subscription -
type Subscription interface {
	Topic() string
	Drain() error
	Delivered() (int64, error)

	Unsubscribe() (int, error)
	AutoUnsubscribe(max int) error

	Pending() (int, int, error)
	MaxPending() (int, int, error)
	PendingLimits() (int, int, error)
	SetPendingLimits(msgLimit, bytesLimit int) error

	// Serve for syncSubscribe
	NextMsg(timeout time.Duration) (Packet, error)
}

// Then I want to us context.Context to manage goroutine

// // Subscription -
// type Subscription interface {
// 	Subscriber

// 	Topics() []string
// }
