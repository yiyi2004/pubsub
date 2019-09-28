package nats

import (
	"context"
	"log"
	"sync"

	pubsub "github.com/DemonDCC/pubsub/interface"
	"github.com/nats.go"
)

var _ pubsub.Broker = &Broker{}

// Broker -
type Broker struct {
	rw   *sync.RWMutex
	URL  string
	Opts *BrokerOptions
	M    map[string]*nats.Conn
}

// NewBroker -
func NewBroker(opts ...nats.Option) *Broker {
	url := NATSURL

	b := &Broker{
		rw: new(sync.RWMutex),
		M:  make(map[string]*nats.Conn),
		Opts: &BrokerOptions{
			Ctx:     context.Background(),
			DefOpts: new(nats.Options),
		},
	}

	if url == "" {
		log.Printf("customed URL is nil, set the URL default URL\n")
		b.URL = nats.DefaultURL
	}

	b.URL = url

	b.Opts.RegisterOptions(opts...)

	return b
}

// CreatePublisher -
func (b *Broker) CreatePublisher(opts ...pubsub.PublisherOptionFunc) pubsub.Publisher {
	p := &Publisher{
		rw: new(sync.RWMutex),
		Opts: &pubsub.PublisherOptions{
			Ctx: context.Background(),
		},
	}

	for _, opt := range opts {
		panic(opt(p.Opts))
	}

	return p
}

// CreateSubscriber -
func (b *Broker) CreateSubscriber(opts ...pubsub.SubscriberOptionFunc) pubsub.Subscriber {
	s := &Subscriber{
		rw: new(sync.RWMutex),
		Opts: &pubsub.SubscriberOptions{
			Ctx: context.Background(),
		},
	}

	for _, opt := range opts {
		panic(opt(s.Opts))
	}

	return s
}

// Topic manager

// Topics -
func (b *Broker) Topics() []string {
	var res []string

	for topic := range b.M {
		res = append(res, topic)
	}

	return res
}

// NumTopics -
func (b *Broker) NumTopics() int {
	return len(b.Topics())
}

// RegisterTopic -
func (b *Broker) RegisterTopic(topic string) (interface{}, error) {
	b.rw.RLock()
	if conn, ok := b.M[topic]; ok {
		b.rw.RUnlock()
		return conn, nil
	}

	conn, err := b.Opts.Connect()
	if err != nil {
		return nil, err
	}

	b.rw.Lock()
	b.M[topic] = conn
	b.rw.Unlock()

	return conn, nil
}

// NumSubcribers -
// If the NumSubscrubers is a negative number
// it represents that connection isn't existed
func (b *Broker) NumSubcribers(topic string) int {
	if conn, ok := b.M[topic]; ok {
		return conn.NumSubscriptions()
	}

	log.Printf("[error]: topic: %s has no corresponding connection\n", topic)
	return -1
}

// Close will close all connections when input len(topics) == 0
// if len(topics) != 0, Close will close the particular connections
func (b *Broker) Close(topics ...string) {
	if len(topics) == 0 {
		for _, conn := range b.M {
			conn.Close()
		}
	}

	for _, topic := range topics {
		if conn, ok := b.M[topic]; ok {
			conn.Close()
		} else {
			log.Printf("[log]: topic %s is not existed\n", topic)
		}
	}
}

// // CreateMultiPublisher -
// func (b *Broker) CreateMultiPublisher(opts ...pubsub.PublisherOptionFunc) pubsub.Publisher {
// 	return &MultiPublisher{
// 		rw:                    new(sync.RWMutex),
// 		DefaultOptionFuncs:    opts,
// 		PublishersOptionFuncs: make(map[string][]pubsub.PublisherOptionFunc),
// 		Publishers:            make(map[string]*Publisher),
// 		Max:                   -1,
// 	}
// }
