package packet

// Packet is a packet in pubsub system
type Packet interface {
	Type() string

	Topic() string
	Payload() []byte

	EncType() string
	Encode(encType string) []byte
}
