package nats

import "github.com/nats.go"

// Msg -
type Msg struct {
	topic string
	reply string
	data  []byte
}

// Payload -
func (msg *Msg) Payload() []byte {
	return msg.data
}

// ReplyT -
func (msg *Msg) ReplyT() string {
	return msg.reply
}

// EncodeMsg -
func EncodeMsg(nmsg *nats.Msg) (*Msg, error) {
	return nil, nil
}
