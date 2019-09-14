package nats

import (
	"errors"
	"log"
	"sync"

	"github.com/nats.go"
)

var (
	errInvalidTopic        = errors.New("[error]: Invalid topic")
	errInvalidChannel      = errors.New("[error]: Invalid channel")
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

	p.rw.RLock()
	if conn, ok := b.M[topic]; ok {
		if err := conn.Publish(topic, data); err != nil {
			p.MsgsNum++

			p.rw.RUnlock()
			return err
		}
	}
	p.rw.RUnlock()

	conn, err := b.Opts.connect()
	if err != nil {
		return err
	}

	p.rw.Lock()
	b.M[topic] = conn
	p.rw.Unlock()

	if err := conn.Publish(topic, data); err == nil {
		if err := conn.Flush(); err != nil {
			return err
		}

		p.MsgsNum++
	}

	return err
}

// PublishMsg -
func (p *Publisher) PublishMsg(b *Broker, msg *natsMsg) error {
	if msg.Subject == "" {
		return errInvalidTopic
	}

	if len(msg.Data) == 0 {
		return errEmptyData
	}

	p.rw.RLock()
	if conn, ok := b.M[msg.Subject]; ok {
		p.rw.RUnlock()

		return conn.PublishMsg((*nats.Msg)(msg))
	}
	p.rw.RUnlock()

	conn, err := b.Opts.connect()
	if err != nil {
		return err
	}

	p.rw.Lock()
	b.M[msg.Subject] = conn
	p.rw.Unlock()

	if err := conn.PublishMsg((*nats.Msg)(msg)); err == nil {
		if err := conn.Flush(); err != nil {
			return err
		}
	}

	return err
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
