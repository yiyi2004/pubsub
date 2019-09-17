package codec

import (
	pkg "github.com/DemonDCC/pubsub/packet"
)

// Encoder -
type Encoder interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}

// Encode -
func Encode(pkg pkg.Packet) ([]byte, error)
