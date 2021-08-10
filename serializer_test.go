package gone

import (
	"fmt"
	"testing"
)

var testValues = []interface{}{
	map[string]string{"foo": "bar"},
	map[string]string{"ping": "pong"},
}

func TestGobSerialization(t *testing.T) {
	var (
		sz           GobEncoder
		serialized   []byte
		deserialized map[string]string
		err          error
	)
	for _, value := range testValues {
		if serialized, err = sz.Serialize(value); err != nil {
			t.Error(err)
		} else {
			deserialized = make(map[string]string)
			if err = sz.Deserialize(serialized, &deserialized); err != nil {
				t.Error(err)
			}
			if fmt.Sprintf("%v", deserialized) != fmt.Sprintf("%v", value) {
				t.Errorf("Expected %v, got %v.", value, deserialized)
			}
		}
	}
}

func TestJSONSerialization(t *testing.T) {
	var (
		sz           JSONEncoder
		serialized   []byte
		deserialized map[string]string
		err          error
	)
	for _, value := range testValues {
		if serialized, err = sz.Serialize(value); err != nil {
			t.Error(err)
		} else {
			deserialized = make(map[string]string)
			if err = sz.Deserialize(serialized, &deserialized); err != nil {
				t.Error(err)
			}
			if fmt.Sprintf("%v", deserialized) != fmt.Sprintf("%v", value) {
				t.Errorf("Expected %v, got %v.", value, deserialized)
			}
		}
	}
}
