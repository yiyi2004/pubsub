package pubsub

// Publisher -
type Publisher interface {
	Publish(b Broker, topic string, data []byte) error
}

// // MultiPublisher -
// type MultiPublisher interface {
// 	Publisher
// }
