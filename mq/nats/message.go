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

// ReplyTopic -
func (msg *Msg) ReplyTopic() string {
	return msg.reply
}
