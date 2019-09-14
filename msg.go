package pubsub

// Msg -
type Msg interface {
	Topic() string
	Payload() []byte
	ReplyTopic() string
	EncoderType() string
	Encode() ([]byte, error)
}
