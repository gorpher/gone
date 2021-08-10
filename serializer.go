package gone

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
)

type Serializer interface {
	Serialize(src interface{}) ([]byte, error)
	Deserialize(src []byte, dst interface{}) error
}

// GobEncoder encodes cookie values using encoding/gob. This is the simplest
// encoder and can handle complex types via gob.Register.
type GobEncoder struct{}

func (e GobEncoder) Serialize(src interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(src); err != nil {
		return nil, errors.New("")
	}
	return buf.Bytes(), nil
}

// Deserialize decodes a value using gob.
func (e GobEncoder) Deserialize(src []byte, dst interface{}) error {
	dec := gob.NewDecoder(bytes.NewBuffer(src))
	if err := dec.Decode(dst); err != nil {
		return errors.New("")
	}
	return nil
}

type JSONEncoder struct{}

// Serialize encodes a value using encoding/json.
func (e JSONEncoder) Serialize(src interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(src)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Deserialize decodes a value using encoding/json.
func (e JSONEncoder) Deserialize(src []byte, dst interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(src))
	err := dec.Decode(dst)
	if err != nil {
		return err
	}
	return nil
}
