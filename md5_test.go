package gone

import (
	"testing"
)

func TestMD5Encrypt(t *testing.T) {
	text := RandBytes(100)
	salt := RandBytes(10)

	actual := MD5Encrypt(text, salt)
	if actual == "" {
		t.Error("MD5Encrypt error")
	}
}
