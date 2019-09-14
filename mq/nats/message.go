package nats

import (
	"github.com/nats.go"
)

// natsMsg -
type natsMsg nats.Msg

// Topic -
func (msg *natsMsg) Topic() string {
	return msg.Subject
}

// Payload -
func (msg *natsMsg) Payload() []byte {
	return msg.Data
}

// ReplyTopic -
func (msg *natsMsg) ReplyTopic() string {
	return msg.Reply
}
