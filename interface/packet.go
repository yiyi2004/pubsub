package pubsub

import "github.com/DemonDCC/pubsub/codec"

// Packet is a packet in pubsub system
type Packet interface {
	// Type is MQ's type; such as nats, kafka, rabbitMQ and so on
	Type() string
	// Topic may be called channel in another Pub/Sub model
	Topic() string
	Payload() []byte
	// if func: Reply is not used, Reply will return empty string
	ReplyTopic() string
	// EncType includes JSON. to be continued
	codec.Encoder
	codec.Decoder
}
