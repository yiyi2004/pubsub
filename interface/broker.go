package pubsub

// Broker -
type Broker interface {
	Topics() []string
	NumTopics() int
	NumSubcribers(topic string) int
	Close() error

	AsyncSubscribe(topic string, handler Handler) (Subscription, error)
	SubscribeSync(topic string, handler Handler) (Subscription, error)

	CreatePublisher(opts ...PublisherOptionFunc) Publisher
	CreateSubscription(opts ...SubscriptionOptionFunc) Subscription
}

// StatusInfo -
type StatusInfo interface {
	ConnStatus() ConnStatus

	// Channel Status

	// Server Status

	// Publisher Status

	// Subscription Status
}

// ConnStatus type
type ConnStatus int8

// Status -
const (
	DisConnected = ConnStatus(iota)
	Connected
	Closed
	Reconnecting
	Connecting

	DrainingSubs
	DrainingPubs
)

// HandlerFunc -
type HandlerFunc func(msg Packet)

// HandlersChain -
type HandlersChain []HandlerFunc

// Handle -
func (h HandlerFunc) Handle(packet Packet) {
	h(packet)
}

// Handler -s
type Handler interface {
	Handle(packet Packet)
}
