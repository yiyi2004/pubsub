package codec

import (
	"bytes"
	"encoding/gob"
)

// Encoder -
type Encoder interface {
	Encode() ([]byte, error)
}

// Decoder -
type Decoder interface {
	Decode(v interface{}) error
}

// GobEncode  -
func GobEncode(v interface{}) ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(v); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GobDecode -
func GobDecode(data []byte, v interface{}) error {
	buf := bytes.NewBuffer(data)

	dec := gob.NewDecoder(buf)

	return dec.Decode(v)
}
