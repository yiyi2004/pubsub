package nats

import (
	"encoding/json"

	"github.com/nats.go"
	pubsub "github.com/zhangce1999/pubsub/interface"
)

var _ pubsub.Packet = &Msg{}

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
	return json.Marshal(m)
}

// Decode -
func (m *Msg) Decode(v interface{}) error {
	return json.Unmarshal(m.Data, v)
}
