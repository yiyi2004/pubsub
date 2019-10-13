package codec

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
)

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

// Encoder -
type Encoder interface {
	Encode() ([]byte, error)
}

// Encode -
func Encode(subj string, enc Encoder) ([]byte, error) {
	return enc.Encode()
}

// Decode -
func Decode(subj string, data []byte, v interface{}) error {
	switch subj {
	case "json":
		return json.Unmarshal(data, v)
	case "gob":
		buf := bytes.NewBuffer(data)

		dec := gob.NewDecoder(buf)

		return dec.Decode(v)
	default:
		return errors.New("[error]: doesn't support this decoding method")
	}
}
