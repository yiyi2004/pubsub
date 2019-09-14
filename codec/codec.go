package codec

// Encoder -
type Encoder interface {
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}
