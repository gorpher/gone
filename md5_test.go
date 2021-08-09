package gone

import (
	"testing"
)

func TestMD5Encrypt(t *testing.T) {
	text := RandBytes(100)
	salt := RandBytes(10)

	actual, err := MD5(text, salt)
	if err != nil {
		t.Fatal(err)
	}
	if actual == "" {
		t.Error("MD5 error")
	}
}

func TestMD5EncryptBytes(t *testing.T) {
	text := RandBytes(100)
	salt := RandBytes(10)

	actual, err := MD5Encrypt(text, salt)
	if err != nil {
		t.Fatal(err)
	}
	if actual == "" {
		t.Error("MD5Encrypt error")
	}
}
