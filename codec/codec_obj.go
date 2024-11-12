package codec

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
)

type ObjectCodec interface {
	codec
	Encode(src interface{}) ([]byte, error)
	Decode(src []byte, dst interface{}) error
}

// GobEncoder encodes cookie values using encoding/gob. This is the simplest
// encoder and can handle complex types via gob.Register.
type GobEncoder struct{}

func (e GobEncoder) Encode(src interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(src); err != nil {
		return nil, errors.New("")
	}
	return buf.Bytes(), nil
}

// Decode decodes a value using gob.
func (e GobEncoder) Decode(src []byte, dst interface{}) error {
	dec := gob.NewDecoder(bytes.NewBuffer(src))
	if err := dec.Decode(dst); err != nil {
		return errors.New("")
	}
	return nil
}

type JSONEncoder struct{}

// Encode encodes a value using encoding/json.
func (e JSONEncoder) Encode(src interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(src)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode decodes a value using encoding/json.
func (e JSONEncoder) Decode(src []byte, dst interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(src))
	err := dec.Decode(dst)
	if err != nil {
		return err
	}
	return nil
}
