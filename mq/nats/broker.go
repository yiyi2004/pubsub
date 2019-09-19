package nats

import (
	"context"
	"log"
	"sync"

	"github.com/DemonDCC/pubsub"
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
func (b *Broker) CreatePublisher(opts ...pubsub.PublisherOptionFunc) *Publisher {
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

// CreateMultiPublisher -
func (b *Broker) CreateMultiPublisher(opts ...pubsub.PublisherOptionFunc) *MultiPublisher {
	return &MultiPublisher{
		rw:                    new(sync.RWMutex),
		DefaultOptionFuncs:    opts,
		PublishersOptionFuncs: make(map[string][]pubsub.PublisherOptionFunc),
		Publishers:            make(map[string]*Publisher),
		Max:                   -1,
	}
}

// CreateSubscriber -
func (b *Broker) CreateSubscriber(opts ...pubsub.SubscriberOptionFunc) *Subscriber {
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

// RegisterTopic -
func (b *Broker) RegisterTopic(topic string) (interface{}, error) {
	if conn, ok := b.M[topic]; ok {
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
