package nats

import (
	"errors"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/zhangce1999/pubsub/codec"
	pubsub "github.com/zhangce1999/pubsub/interface"
)

var (
	errInvalidTopic        = errors.New("[error]: invalid topic")
	errInvalidChannel      = errors.New("[error]: invalid channel")
	errInvalidConnection   = errors.New("[error]: invalid connection")
	errInvalidBroker       = errors.New("[error]: invalid broker")
	errInvalidPublisher    = errors.New("[error]: invalid publisher")
	errInvalidSubscriber   = errors.New("[error]: invalid subscriber")
	errInvalidMultiPublish = errors.New("[error]: multiPublish error")
	errEmptyData           = errors.New("[error]: empty data")
)

var _ pubsub.Publisher = &Publisher{}

// Publisher -
type Publisher struct {
	// rw represents a Read/Write Mutex
	rw      *sync.Mutex
	Topic   string
	MsgsNum int
	Opts    *pubsub.PublisherOptions
}

// MultiPublisher -
type MultiPublisher struct {
	rw *sync.Mutex

	DefaultOptionFuncs []pubsub.PublisherOptionFunc

	PublishersOptionFuncs map[string][]pubsub.PublisherOptionFunc
	Publishers            map[string]*Publisher

	// Max represents the maximum value of publisher
	Max int
}

// Publish will publish raw data to topic
func (p *Publisher) Publish(b pubsub.Broker, topic string, data []byte) error {
	log.Printf("[pub]: publish message to topic: %s\n", topic)
	broker, ok := b.(*Broker)
	if !ok {
		return errInvalidBroker
	}

	if topic != "" {
		p.Topic = topic
	} else {
		return errInvalidTopic
	}

	if len(data) == 0 {
		return errEmptyData
	}

	return p.publish(broker, data)
}

func (p *Publisher) publish(b *Broker, data []byte) error {
	conn, err := b.RegisterTopic(p.Topic)
	if err != nil {
		return err
	}

	if max, ok := p.Opts.Ctx.Value("MAX_MESSAGES").(int); ok {
		if p.MsgsNum >= max && max > 0 {
			log.Fatal(`[log]: the amount of messages that can be sent have
						have reached the maximum`)
		}
	}

	if conn, ok := conn.(*nats.Conn); ok {
		encData, err := encode(p.Topic, data)
		if err != nil {
			return err
		}

		if err := conn.Publish(p.Topic, encData); err == nil {
			if err := conn.Flush(); err != nil {
				return err
			}

			p.MsgsNum++

			return nil
		}
	}

	return errInvalidConnection
}

// PublishMsg will be abondoned
func (p *Publisher) PublishMsg(b pubsub.Broker, pkg pubsub.Packet) error {
	broker, ok := b.(*Broker)
	if !ok {
		return errInvalidBroker
	}

	if pkg.Topic() == "" {
		return errInvalidTopic
	}

	if len(pkg.Payload()) == 0 {
		return errEmptyData
	}

	return p.publishMsg(broker, pkg)
}

func (p *Publisher) publishMsg(b *Broker, pkg pubsub.Packet) error {
	conn, err := b.RegisterTopic(pkg.Topic())
	if err != nil {
		return err
	}

	if conn, ok := conn.(nats.Conn); ok {
		if err := conn.Publish(pkg.Topic(), pkg.Payload()); err == nil {
			if err := conn.Flush(); err != nil {
				return err
			}

			p.MsgsNum++

			return nil
		}
	}

	return errInvalidConnection
}

// MultiPublish -
func (mp *MultiPublisher) MultiPublish(b pubsub.Broker, topics []string, data [][]byte) error {
	if len(topics) != len(data) {
		return errInvalidMultiPublish
	}

	for i := 0; i < len(topics); i++ {
		if topics[i] == "" {
			log.Println("invalid topic")
			continue
		}

		topic := topics[i]

		if p, ok := mp.Publishers[topic]; ok {
			go p.Publish(b, topic, data[i])
		} else {
			if opts, ok := mp.PublishersOptionFuncs[topic]; ok {
				p := b.CreatePublisher(opts...)
				publisher, ok := p.(*Publisher)
				if !ok {
					return errInvalidPublisher
				}

				mp.rw.Lock()
				mp.Publishers[topic] = publisher
				mp.rw.Unlock()

				go p.Publish(b, topic, data[i])
			}
		}
	}

	return nil
}

func encode(topic string, data []byte) (res []byte, err error) {
	msg := Msg{}
	msg.Subject = topic
	msg.Data = data

	return codec.GobEncode(&msg)
}
