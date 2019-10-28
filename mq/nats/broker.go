package nats

import (
	"context"
	"sync"

	"github.com/nats-io/nats.go"
	pubsub "github.com/zhangce1999/pubsub/interface"
)

var _ pubsub.Broker = &Broker{}
var _ pubsub.Group = &Broker{}

// Broker -
type Broker struct {
	URL             string
	Opts            *BrokerOptions
	DefaultHandlers pubsub.HandlersChain
	Group

	topics []string
	tree   Trie
	rw     *sync.Mutex
	sch    chan *nats.Msg
}

// NewBroker -
func NewBroker(opts ...nats.Option) *Broker {
	b := &Broker{
		rw: new(sync.Mutex),
		Opts: &BrokerOptions{
			Ctx:     context.Background(),
			DefOpts: new(nats.Options),
		},
		sch:  make(chan *nats.Msg),
		tree: NewTrie('/'),
		Group: Group{
			root:     true,
			basePath: "/",
			Handlers: nil,
		},
	}

	b.Group.broker = b

	if NATSURL == "" {
		b.URL = nats.DefaultURL
	} else {
		b.URL = NATSURL
	}

	b.Opts.RegisterOptions(opts...)

	return b
}

// CreatePublisher -
func (b *Broker) CreatePublisher(opts ...pubsub.PublisherOptionFunc) pubsub.Publisher {
	p := &Publisher{
		rw: new(sync.Mutex),
		Opts: &pubsub.PublisherOptions{
			Ctx: context.Background(),
		},
	}

	for _, opt := range opts {
		panic(opt(p.Opts))
	}

	if topic, ok := p.Opts.Ctx.Value(pubsub.Key("Topic")).(string); ok {
		p.Topic = topic
	}

	return p
}

// CreateSubscription -
func (b *Broker) CreateSubscription(opts ...pubsub.SubscriptionOptionFunc) pubsub.Subscription {
	s := &Subscription{
		rw: new(sync.Mutex),
		Opts: &pubsub.SubscriberOptions{
			Ctx: context.Background(),
		},
		Subs: make(map[string]*nats.Subscription),
	}

	for _, opt := range opts {
		panic(opt(s.Opts))
	}

	if topic, ok := s.Opts.Ctx.Value(pubsub.Key("Topic")).(string); ok {
		s.topic = topic
	}
	if isGroup, ok := s.Opts.Ctx.Value(pubsub.Key("IsGroup")).(bool); ok {
		s.isGroup = isGroup
	}
	if subType, ok := s.Opts.Ctx.Value(pubsub.Key("SubscriptionType")).(pubsub.SubscriptionType); ok {
		s.subtype = subType
	}

	return s
}

// Topics -
func (b *Broker) Topics() []string {
	b.rw.Lock()
	defer b.rw.Unlock()
	return b.topics
}

// NumTopics -
func (b *Broker) NumTopics() int {
	b.rw.Lock()
	defer b.rw.Unlock()
	return len(b.topics)
}

// NumSubcribers -
func (b *Broker) NumSubcribers(topic string) int {
	return 0
}

// Close -
func (b *Broker) Close() error

// AsyncSubscribe -
func (b *Broker) AsyncSubscribe(ctx context.Context, topic string, handler pubsub.HandlerFunc) (pubsub.Subscription, error)

// SubscribeSync -
func (b *Broker) SubscribeSync(ctx context.Context, topic string, handler pubsub.HandlerFunc) (pubsub.Subscription, error)

func (b *Broker) QueueSubscribeSync(ctx context.Context, topic string, queue string) (Subscription, error)
