package nats

import (
	"github.com/DemonDCC/pubsub/packet"
	"github.com/nats.go"
)

var _ packet.Packet = &Msg{}

// Msg -
type Msg struct {
	nats.Msg
}

// Type -
func (m *Msg) Type() string {
	return "nats"
}

// Topic -
func (m *Msg) Topic() string {
	return m.Subject
}

// Payload -
func (m *Msg) Payload() []byte {
	return m.Data
}

// ReplyTopic -
func (m *Msg) ReplyTopic() string {
	return m.Reply
}

// Encode -
func (m *Msg) Encode() ([]byte, error) {
	return nil, nil
}

// Decode -
func (m *Msg) Decode(v interface{}) error {
	return nil
}
