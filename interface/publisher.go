package pubsub

// Publisher -
type Publisher interface {
	Publish(b Broker, topic string, data []byte) error
	PublishMsg(b Broker, packet Packet) error
}

// // MultiPublisher -
// type MultiPublisher interface {
// 	Publisher
// }
