package nats

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/zhangce1999/pubsub/codec"
	pubsub "github.com/zhangce1999/pubsub/interface"
)

var (
	errMultiSubscribe = errors.New("[error]: MultiSubscribe error")
)

var _ pubsub.Subscription = &Subscription{}

// Subscription -
type Subscription struct {
	isGroup bool
	rw      *sync.Mutex
	topic   string
	broker  *Broker
	// Sub is a nats.Subscription which represensts interest in the given topic.
	Subs []*nats.Subscription
	Opts *pubsub.SubscriberOptions
}

// Topic -
func (s *Subscription) Topic() string {
	return s.topic
}

// Drain -
func (s *Subscription) Drain() error {
	if s.isGroup {
		return s.Subs[0].Drain()
	}

	if routes, ok := s.broker.M[s.topic]; ok {
		for _, route := range routes {
			go func() {
				if err := route.conn.Drain(); err != nil {
					log.Println(err)
				}
			}()
		}
	}

	return
}

// Delivered -
func (s *Subscription) Delivered() (int64, error) {
	return s.Sub.Delivered()
}

// Unsubscribe -
func (s *Subscription) Unsubscribe() (int, error) {
	return 1, s.Sub.Unsubscribe()
}

// AutoUnsubscribe -
func (s *Subscription) AutoUnsubscribe(max int) error {
	return s.AutoUnsubscribe(max)
}

// Pending -
func (s *Subscription) Pending() (int, int, error) {
	return s.Sub.Pending()
}

// MaxPending -
func (s *Subscription) MaxPending() (int, int, error) {
	return s.Sub.MaxPending()
}

// PendingLimits -
func (s *Subscription) PendingLimits() (int, int, error) {
	return s.Sub.PendingLimits()
}

// SetPendingLimits -
func (s *Subscription) SetPendingLimits(msgLimit, bytesLimit int) error {
	return s.Sub.SetPendingLimits(msgLimit, bytesLimit)
}

// NextMsg -
func (s *Subscription) NextMsg(timeout time.Duration) (res pubsub.Packet, err error) {
	msg, err := s.Sub.NextMsg(timeout)
	if err != nil {
		return nil, err
	}

	if err := codec.GobDecode(msg.Data, &res); err != nil {
		return nil, err
	}

	return
}

// // ChanSubscribe -
// func (s *Subscription) ChanSubscribe(b pubsub.Broker, topic string, ch interface{}) error {
// 	log.Printf("[sub]: subscribe topic: %s\n", topic)
// 	broker, ok := b.(*Broker)
// 	if !ok {
// 		return errInvalidSubscriber
// 	}

// 	channel, ok := ch.(chan *nats.Msg)
// 	if !ok {
// 		return errInvalidChannel
// 	}

// 	if topic != "" {
// 		s.topic = topic
// 	} else {
// 		return errInvalidTopic
// 	}

// 	if ch == nil {
// 		return errInvalidChannel
// 	}

// 	s.rw.Lock()
// 	conn, err := broker.RegisterTopic(topic)
// 	if err != nil {
// 		return err
// 	}
// 	s.rw.Unlock()

// 	if conn, ok := conn.(*nats.Conn); ok {
// 		sub, err := conn.ChanSubscribe(topic, channel)
// 		if err != nil {
// 			return err
// 		}

// 		s.Sub = sub
// 	}

// 	return errInvalidConnection
// }

// // Unsubscribe -
// func (s *Subscription) Unsubscribe() (int, error) {
// 	return 1, s.Sub.Unsubscribe()
// }

// // AutoUnsubscribe -
// func (s *Subscription) AutoUnsubscribe(max int) error {
// 	return s.Sub.AutoUnsubscribe(max)
// }

// // MultiSubscriber is either a queue (with the same topic) or a set of Subscribers(with different topics)
// type MultiSubscriber struct {
// 	wg *sync.WaitGroup

// 	DefaultOptionFuncs []pubsub.SubscriberOptionFunc

// 	SubscriberOptionFuncs map[string][]pubsub.SubscriberOptionFunc
// 	Subscribers           map[string]*Subscription

// 	NumSubs int

// 	// isQueue   bool
// 	// QueueName string // queue

// 	// Only when isQueue=true, Topic is set.
// 	Topic string

// 	// Max represents the maximum number of Subscribers
// 	// if Max is a negative number, it represents that
// 	// the number of Subscribers is unlimited.
// 	Max int
// }

// // ChanSubscribe : if the value of parameter queue is "", it represents that
// // the MultiSubscriber is a set of Subscribers, otherwise the MultiSubscriber is
// // Queue of Subscribers
// func (ms *MultiSubscriber) ChanSubscribe(b *Broker, topic string, ch chan *nats.Msg) error {
// 	if topic == "" {
// 		return errInvalidTopic
// 	}

// 	if ch == nil {
// 		return errInvalidChannel
// 	}

// 	return ms.chanSubscribe(b, topic, ch)
// }

// // MultiChanSubscribe -
// func (ms *MultiSubscriber) MultiChanSubscribe(b *Broker, topics []string, chs []chan *nats.Msg) error {
// 	if len(topics) != len(chs) {
// 		return errMultiSubscribe
// 	}

// 	for i, topic := range topics {
// 		if topic == "" {
// 			log.Printf("[log]: topics[%d] is empty string\n", i)
// 			continue
// 		}

// 		if chs[i] == nil {
// 			log.Printf("[log]: channel[%d] is nil", i)
// 			continue
// 		}

// 		if s, ok := ms.Subscribers[topic]; ok {
// 			if s.Sub.IsValid() {
// 				log.Printf("[log]: topic %s has been subscribe", topic)
// 				continue
// 			}
// 		}

// 		if opts, ok := ms.SubscriberOptionFuncs[topic]; ok {
// 			s := b.CreateSubscriber(opts...)

// 			s.rw.Lock()
// 			ms.Subscribers[topic] = s
// 			s.rw.Unlock()

// 			if err := s.ChanSubscribe(b, topic, chs[i]); err != nil {
// 				log.Printf("[error]: Subscribe %s failed", topic)
// 			}
// 		}
// 	}

// 	return nil
// }

// func (ms *MultiSubscriber) chanSubscribe(b pubsub.Broker, topic string, ch chan *nats.Msg) error {
// 	if s, ok := ms.Subscribers[topic]; ok {
// 		return s.ChanSubscribe(b, topic, ch)
// 	}

// 	if opts, ok := ms.SubscriberOptionFuncs[topic]; ok {
// 		sub := b.CreateSubscriber(opts...)

// 		sub.rw.Lock()
// 		ms.Subscribers[topic] = sub
// 		sub.rw.Unlock()

// 		return sub.ChanSubscribe(b, topic, ch)
// 	}

// 	sub := b.CreateSubscriber(ms.DefaultOptionFuncs...)

// 	sub.rw.Lock()
// 	ms.Subscribers[topic] = sub
// 	sub.rw.Unlock()

// 	return sub.ChanSubscribe(b, topic, ch)
// }

// // Wait -
// func (ms *MultiSubscriber) Wait() {
// 	ms.wg.Wait()
// }
