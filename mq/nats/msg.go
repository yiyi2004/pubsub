package nats

import (
	pubsub "github.com/zhangce1999/pubsub/interface"
)

var _ pubsub.Packet = &Msg{}

// Msg -
type Msg struct {
	topic string
	reply string
	data  []byte
}

// Type -
func (m *Msg) Type() string {
	return "nats"
}

// Topic -
func (m *Msg) Topic() string {
	return m.topic
}

// Payload -
func (m *Msg) Payload() []byte {
	return m.data
}

// ReplyTopic -
func (m *Msg) ReplyTopic() string {
	return m.reply
}

// Encode -
func (m *Msg) Encode() ([]byte, error) {
	return nil, nil
}

// Decode -
func (m *Msg) Decode(v interface{}) error {
	return nil
}
