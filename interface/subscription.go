package pubsub

import "time"

// SubscriptionType -
type SubscriptionType int

// SubscriptionType -
const (
	SyncSubscription = SubscriptionType(iota)
	AsyncSubscription
)

// Subscription -
type Subscription interface {
	Type() SubscriptionType
	Topics() []string

	Unsubscribe(topics ...string) (int, error)
	AutoUnsubscribe(max int, topic string) error

	// Serve for SubscribeSync
	NextMsg(timeout time.Duration, topic string) (Packet, error)
}
