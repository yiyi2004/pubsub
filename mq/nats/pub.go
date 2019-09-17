package nats

import (
	"errors"
	"log"
	"sync"

	pkg "github.com/DemonDCC/pubsub/packet"
	"github.com/nats.go"
)

var (
	errInvalidTopic        = errors.New("[error]: Invalid topic")
	errInvalidChannel      = errors.New("[error]: Invalid channel")
	errInvalidConnection   = errors.New("[error]: Invalid Connection")
	errInvalidMultiPublish = errors.New("[error]: MultiPublish error")
	errEmptyData           = errors.New("[error]: Empty Data")
)

// Publisher -
type Publisher struct {
	// rw represents a Read/Write Mutex
	rw      sync.RWMutex
	Topic   string
	MsgsNum int
	Opts    *PublisherOptions
}

// MultiPublisher -
type MultiPublisher struct {
	rw sync.RWMutex

	DefaultOptionFuncs []PublisherOptionFunc

	PublishersOptionFuncs map[string][]PublisherOptionFunc
	Publishers            map[string]*Publisher

	// Max represents the maximum value of publisher
	Max int
}

// Publish will publish raw data to topic
func (p *Publisher) Publish(b *Broker, topic string, data []byte) error {
	if topic != "" {
		p.Topic = topic
	} else {
		return errInvalidTopic
	}

	if len(data) == 0 {
		return errEmptyData
	}

	return p.publish(b, topic, data)
}

func (p *Publisher) publish(b *Broker, topic string, data []byte) error {
	conn, err := b.RegisterTopic(topic)
	if err != nil {
		return err
	}

	if conn, ok := conn.(*nats.Conn); ok {
		if err := conn.Publish(topic, data); err == nil {
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
func (p *Publisher) PublishMsg(b *Broker, pkg pkg.Packet) error {
	if pkg.Topic() == "" {
		return errInvalidTopic
	}

	if len(pkg.Payload()) == 0 {
		return errEmptyData
	}

	return p.publishMsg(b, pkg)
}

func (p *Publisher) publishMsg(b *Broker, pkg pkg.Packet) error {
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
func (mp *MultiPublisher) MultiPublish(b *Broker, topics []string, data [][]byte) error {
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

				mp.rw.Lock()
				mp.Publishers[topic] = p
				mp.rw.Unlock()

				go p.Publish(b, topic, data[i])
			}
		}
	}

	return nil
}
