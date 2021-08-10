package gone

import (
	"testing"
)

func TestMD5Encrypt(t *testing.T) {
	text := RandBytes(100)
	salt := RandBytes(10)
	actual := MD5(text, salt)
	if actual == "" {
		t.Error("MD5 error")
	}
}

func TestMD5EncryptBytes(t *testing.T) {
	text := RandBytes(100)
	salt := RandBytes(10)

	actual := HMacMD5(text, salt)

	if actual == "" {
		t.Error("HMacMD5 error")
	}
}

func TestHMacSha256(t *testing.T) {
	text := RandBytes(100)
	salt := RandBytes(10)
	actual := HMacSha256(text, salt)
	if len(actual) == 0 {
		t.Error("HMacMD5 error")
	}
}
