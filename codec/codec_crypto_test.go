package codec

import (
	"bytes"
	"encoding/json"
	"github.com/gorpher/gone/core"
	"github.com/gorpher/gone/osutil"
	"reflect"
	"testing"
	"time"
)

func TestCookieCodec(t *testing.T) {
	var tests []CryptoCodec
	for i := 0; i < 10; i++ {
		tests = append(tests, NewCookieCodec([]byte(osutil.RandString(8)), osutil.RandBytes(16)))
	}
	for i := range tests {
		s := tests[i]
		value := map[string]interface{}{
			osutil.RandString(4): osutil.RandString(20),
			osutil.RandString(5): osutil.RandString(12),
			osutil.RandString(6): osutil.RandString(8),
		}
		encoder := JSONEncoder{}
		data, err := encoder.Encode(value)
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
		err = encoder.Decode(plaintData, &dst)
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
		ExpirationTime: &core.Time{time.Unix(1628603180, 0)},
		NotBefore:      &core.Time{time.Unix(1628603180, 0)},
		IssuedAt:       &core.Time{time.Unix(1628603180, 0)},
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
