package nats

import (
	"errors"
	"log"
	"sync"

	"github.com/nats.go"
	pubsub "github.com/zhangce1999/pubsub/interface"
)

var (
	errMultiSubscribe = errors.New("[error]: MultiSubscribe error")
)

var _ pubsub.Subscriber = &Subscriber{}

// Subscriber -
type Subscriber struct {
	rw    *sync.Mutex
	topic string
	// Sub is a nats.Subscription which represensts interest in the given topic.
	Sub  *nats.Subscription
	Opts *pubsub.SubscriberOptions
}

// Topic -
func (s *Subscriber) Topic() string {
	return s.topic
}

// ChanSubscribe -
func (s *Subscriber) ChanSubscribe(b pubsub.Broker, topic string, ch interface{}) error {
	log.Printf("[sub]: subscribe topic: %s\n", topic)
	broker, ok := b.(*Broker)
	if !ok {
		return errInvalidSubscriber
	}

	channel, ok := ch.(chan *nats.Msg)
	if !ok {
		return errInvalidChannel
	}

	if topic != "" {
		s.topic = topic
	} else {
		return errInvalidTopic
	}

	if ch == nil {
		return errInvalidChannel
	}

	s.rw.Lock()
	conn, err := broker.RegisterTopic(topic)
	if err != nil {
		return err
	}
	s.rw.Unlock()

	if conn, ok := conn.(*nats.Conn); ok {
		sub, err := conn.ChanSubscribe(topic, channel)
		if err != nil {
			return err
		}

		s.Sub = sub
	}

	return errInvalidConnection
}

// Unsubscribe -
func (s *Subscriber) Unsubscribe() (int, error) {
	return 1, s.Sub.Unsubscribe()
}

// AutoUnsubscribe -
func (s *Subscriber) AutoUnsubscribe(max int) error {
	return s.Sub.AutoUnsubscribe(max)
}

// MultiSubscriber is either a queue (with the same topic) or a set of Subscribers(with different topics)
type MultiSubscriber struct {
	wg *sync.WaitGroup

	DefaultOptionFuncs []pubsub.SubscriberOptionFunc

	SubscriberOptionFuncs map[string][]pubsub.SubscriberOptionFunc
	Subscribers           map[string]*Subscriber

	NumSubs int

	// isQueue   bool
	// QueueName string // queue

	// Only when isQueue=true, Topic is set.
	Topic string

	// Max represents the maximum number of Subscribers
	// if Max is a negative number, it represents that
	// the number of Subscribers is unlimited.
	Max int
}

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
