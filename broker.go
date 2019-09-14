package pubsub

// Broker -
type Broker interface {
	Topic

	CreatePublisher(options ...Option)
}

// Topic -
type Topic interface {
	Topics() []string
	RegisterTopic(topic string) error
	Topic(topic string) (conn interface{}, ok bool)
	NumTopics() int
	NumSubcribers(topic string) int
}

// Metrics , emmm I don't konw how to explain the Metrics interface
type Metrics interface {
}

// ConnManager is an interaface to manage connection, one connection is related to a topic
type ConnManager interface {
	Reconnect(topic string) error

	Shutdown(topic string) error
	ShutdownAll() (int, error)
}

// Status -
type Status interface {
	// Connection statics

	// Channel statics

	// Server statics

	// Publisher statics

	// Subscription statics
}
