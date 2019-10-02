package codec

// Encoder -
type Encoder interface {
	Encode() ([]byte, error)
}

// Decoder -
type Decoder interface {
	Decode(v interface{}) error
}
