package nats

import (
	"context"
	"sync"
	"time"

	pubsub "github.com/zhangce1999/pubsub/interface"
)

// type Subscription interface {
// 	pubsub.Subscription

// }

var _ pubsub.Subscription = &Subscription{}

// Subscription -
type Subscription struct {
	// Sub is a nats.Subscription which represensts interest in the given topic.
	Opts *pubsub.SubscriberOptions

	group
	rw     *sync.Mutex
	broker *Broker

	tree    Trie
	subtype pubsub.SubscriptionType
}

// Type -
func (s *Subscription) Type() pubsub.SubscriptionType

// Topics -
func (s *Subscription) Topics() []string

// Unsubscribe -
func (s *Subscription) Unsubscribe(topics ...string) (int, error)

// AutoUnsubscribe -
func (s *Subscription) AutoUnsubscribe(max int, topic string) error

// Close -
func (s *Subscription) Close()

// Filter -
func (s *Subscription) Filter(ctx context.Context, in chan pubsub.Packet, quit chan struct{}, filters ...func(pubsub.Packet) bool) (out chan pubsub.Packet)

// NextMsg -
func (s *Subscription) NextMsg(timeout time.Duration, topic string, out chan pubsub.Packet, errChan chan error)
