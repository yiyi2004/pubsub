package packet

import "github.com/nats.go"

// Packet represents a packet can contain any type of message
// Packet should have a type field to record the message type
// and another type to record the packet type such as: ...
// emmm...
type Packet interface {
	PkgName() string
	Type() string
	Payload() []byte
}

// NatsPacket -
type NatsPacket struct {
	Name string
	Type string

	Msg *nats.Msg
}
