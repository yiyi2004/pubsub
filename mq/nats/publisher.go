package nats

import (
	"context"
	"errors"
	"sync"

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

	cancel context.CancelFunc
}

// Publish -
func (p *Publisher) Publish(ctx context.Context, b pubsub.Broker, topic string, in chan string, errChan chan error)

// PublishRequest -
func (p *Publisher) PublishRequest(ctx context.Context, b pubsub.Broker, topic string, reply string, in chan string, errChan chan error)

// Flush -
func (p *Publisher) Flush() error

// Close -
func (p *Publisher) Close()
