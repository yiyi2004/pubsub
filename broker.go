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

// ConnManager is an interaface to manage connection, one connection is related to a topic
type ConnManager interface {
	Reconnect(topic string) error

	Shutdown(topic string) error
	ShutdownAll() (int, error)
}

// Status -
type Status interface {
	ConnStatus() *ConnStatus

	// Channel Status

	// Server Status

	// Publisher Status

	// Subscription Status
}

// ConnStatus -
type ConnStatus struct {
	DisConnected bool
	Connected    bool
	Closed       bool
	Reconnecting bool
	Connecting   bool

	DrainingSubs bool
	DrainingPubs bool
}
