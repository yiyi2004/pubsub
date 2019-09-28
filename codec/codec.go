package codec

import (
	pkg "github.com/DemonDCC/pubsub/packet"
)

// Encoder -
type Encoder interface {
	Encode() ([]byte, error)
}

// Decoder -
type Decoder interface {
	Decode(v interface{}) error
}

// Encode -
func Encode(pkg pkg.Packet) ([]byte, error)
