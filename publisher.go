package pubsub

import "context"

// Publisher -
type Publisher interface {
	Publish(topic string, data []byte) error
	PublishMsg(msg Msg) error

	PublishWithContext(ctx context.Context, topic string, data []byte) error
	PublishMsgWithContext(ctx context.Context, msg Msg) error
}

// MultiPublisher -
type MultiPublisher interface {
	Publisher
}
