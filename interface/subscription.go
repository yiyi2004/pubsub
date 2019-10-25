package pubsub

import (
	"context"
	"time"
)

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
	Close()

	Filter(ctx context.Context, in chan Packet, quit chan struct{}, filters ...func(Packet) bool) (out chan Packet)
	// Serve for SubscribeSync and AsyncSusbcribe
	NextMsg(timeout time.Duration, topic string, out chan Packet, errChan chan error)
}
