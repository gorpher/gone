package gone

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestCookieCodec(t *testing.T) {
	var tests []Codec
	for i := 0; i < 10; i++ {
		tests = append(tests, NewCookieCodec([]byte(RandString(8)), RandBytes(16)))
	}
	for i := range tests {
		s := tests[i]
		value := map[string]interface{}{
			RandString(4): RandString(20),
			RandString(5): RandString(12),
			RandString(6): RandString(8),
		}
		encoder := JSONEncoder{}
		data, err := encoder.Serialize(value)
		if err != nil {
			t.Fatal(err)
		}
		encoded, err1 := s.Encode([]byte("sid"), data)
		if err1 != nil {
			t.Fatal(err1)
		}
		plaintData, err2 := s.Decode([]byte("sid"), encoded)
		if err2 != nil {
			t.Fatalf("%v: %v", err2, encoded)
		}
		if !bytes.Equal(plaintData, data) {
			t.Fatal("encode error")
		}
		dst := make(map[string]interface{})
		err = encoder.Deserialize(plaintData, &dst)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(value, dst) {
			t.Fatalf("Expected %v, got %v.", value, dst)
		}
	}
}

func TestJwtCodec(t *testing.T) {
	codec := NewJwtCodec("HS256")
	v, err := json.Marshal(Payload{
		Subject:        "access_token",
		Issuer:         "gorpher",
		ExpirationTime: &Time{time.Unix(1628603180, 0)},
		NotBefore:      &Time{time.Unix(1628603180, 0)},
		IssuedAt:       &Time{time.Unix(1628603180, 0)},
		Audience:       Audience{"https://www.gorpher.site/"},
	})
	if err != nil {
		t.Fatal(err)
	}
	key := []byte("123456")
	t.Log(string(key))
	encode, err := codec.Encode(key, v)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(encode))

	decode, err := codec.Decode(key, encode)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(decode))
}

func equalMap(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v1 := range a {
		v2, ok := b[k]
		if !ok {
			return false
		}
		if !reflect.DeepEqual(v1, v2) {
			return false
		}
	}
	return true
}
