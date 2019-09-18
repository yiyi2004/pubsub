package pubsub

import (
	"context"

	pkg "github.com/DemonDCC/pubsub/packet"
)

// Publisher -
type Publisher interface {
	Publish(b Broker, topic string, data []byte) error
	PublishMsg(b Broker, packet pkg.Packet) error

	PublishWithContext(ctx context.Context, b Broker, topic string, data []byte) error
	PublishMsgWithContext(ctx context.Context, b Broker, packet pkg.Packet) error
}

// MultiPublisher -
type MultiPublisher interface {
	Publisher
}
