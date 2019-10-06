package nats

import (
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/zhangce1999/pubsub/codec"
	pubsub "github.com/zhangce1999/pubsub/interface"
)

// type Subscription interface {
// 	pubsub.Subscription

// }

var _ pubsub.Subscription = &Subscription{}

// Subscription -
type Subscription struct {
	isGroup bool
	rw      *sync.Mutex
	topic   string
	broker  *Broker
	subtype pubsub.SubscriptionType
	// Sub is a nats.Subscription which represensts interest in the given topic.
	Subs map[string]*nats.Subscription
	Opts *pubsub.SubscriberOptions
}

// Type -
func (s *Subscription) Type() pubsub.SubscriptionType {
	return s.subtype
}

// Topics -
func (s *Subscription) Topics() []string {
	s.rw.Lock()
	defer s.rw.Unlock()

	var topics []string

	if routes, ok := s.broker.M[s.topic]; ok {
		for relativePath := range routes {
			topics = append(topics, joinPaths(s.topic, relativePath))
		}
		return topics
	}

	topics = append(topics, s.topic)
	return topics
}

// Unsubscribe -
func (s *Subscription) Unsubscribe(topics ...string) (num int, err error) {
	s.rw.Lock()
	defer s.rw.Unlock()

	if s.isGroup {
		if len(topics) == 0 {
			for relativePath, sub := range s.Subs {
				num++

				sub.Unsubscribe()
				s.broker.M[s.topic][relativePath].conn.Close()

				delete(s.Subs, relativePath)
				delete(s.broker.M[s.topic], relativePath)
			}

			return num, nil
		}

		for _, topic := range topics {
			topic = insertBackSlash(topic)
			for relativePath, sub := range s.Subs {
				targetTopic := joinPaths(s.topic, relativePath)
				if topic == targetTopic {
					num++

					sub.Unsubscribe()
					s.broker.M[s.topic][relativePath].conn.Close()

					delete(s.Subs, relativePath)
					delete(s.broker.M[s.topic], relativePath)
				}
			}
		}

		return num, nil
	}

	if err = s.Subs[s.topic].Unsubscribe(); err != nil {
		return
	}

	delete(s.Subs, s.topic)
	delete(s.broker.M["/"], s.topic)
	return
}

// AutoUnsubscribe -
func (s *Subscription) AutoUnsubscribe(max int, topic string) error {
	s.rw.Lock()
	defer s.rw.Unlock()

	topic = insertBackSlash(topic)
	if s.isGroup {
		if s.topic == topic {
			for _, sub := range s.Subs {
				sub.AutoUnsubscribe(max)
			}

			s.Subs = make(map[string]*nats.Subscription)
		}

		for relativePath, sub := range s.Subs {
			targetTopic := joinPaths(s.topic, relativePath)
			if topic == targetTopic {
				return sub.AutoUnsubscribe(max)
			}
		}
	}

	return s.Subs[s.topic].AutoUnsubscribe(max)
}

// NextMsg -
func (s *Subscription) NextMsg(timeout time.Duration, topic string) (pubsub.Packet, error) {
	s.rw.Lock()
	defer s.rw.Unlock()

	if !s.isGroup {
		msg, err := s.Subs[s.topic].NextMsg(timeout)
		if err != nil {
			return nil, err
		}

		var res Msg

		if err := codec.GobDecode(msg.Data, &res); err != nil {
			return nil, err
		}

		return &res, nil
	}

	topic = insertBackSlash(topic)
	for relativePath, sub := range s.Subs {
		targetTopic := joinPaths(s.topic, relativePath)
		if topic == targetTopic {
			msg, err := sub.NextMsg(timeout)
			if err != nil {
				return nil, err
			}

			var res Msg

			if err := codec.GobDecode(msg.Data, &res); err != nil {
				return nil, err
			}

			return &res, nil
		}
	}

	return nil, errInvalidTopic
}
