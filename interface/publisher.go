package pubsub

import "context"

// Publisher -
type Publisher interface {
	Publish(ctx context.Context, b Broker, topic string, in chan string, errChan chan error)
	PublishRequest(ctx context.Context, b Broker, topic string, reply string, in chan string, errChan chan error)
	Flush() error
	Close()
}
