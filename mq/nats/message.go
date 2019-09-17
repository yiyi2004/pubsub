package nats

// Msg -
type Msg struct {
	topic   string
	reply   string
	payload []byte
}

// Payload -
func (msg *Msg) Payload() []byte {
	return msg.payload
}

// ReplyT -
func (msg *Msg) ReplyT() string {
	return msg.reply
}
