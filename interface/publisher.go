package pubsub

import (
	pkg "github.com/DemonDCC/pubsub/packet"
)

// Publisher -
type Publisher interface {
	Publish(b Broker, topic string, data []byte) error
	PublishMsg(b Broker, packet pkg.Packet) error
}

// MultiPublisher -
type MultiPublisher interface {
	Publisher
}
