package pubsub

// Packet is a packet in pubsub system
type Packet interface {
	// Topic may be called channel in another Pub/Sub model
	Topic() string
	Payload() []byte
	// if func: ReplyTopic is not used, ReplyTopic will return empty string
	ReplyTopic() string
}
