package core_test

import (
	bytes2 "bytes"
	"github.com/gorpher/gone/core"
	"testing"
)

var licenseStr = "需要签名的字符串"

func TestBase64(t *testing.T) {
	bytes := []byte(licenseStr)
	encode := core.Base64StdEncode(bytes)
	decode, err := core.Base64StdDecode(encode)
	if err != nil {
		t.Fatal("Base64StdDecode error")
	}
	if !bytes2.Equal(decode, bytes) {
		t.Fatalf("Base64StdEncode error: want is %s, but acutal is %s", bytes, decode)
	}

	encode = core.Base64RawStdEncode(bytes)
	decode, err = core.Base64RawStdDecode(encode)
	if err != nil {
		t.Fatal("Base64RawStdDecode error")
	}
	if !bytes2.Equal(decode, bytes) {
		t.Fatalf("Base64RawStdDecode error: want is %s, but acutal is %s", bytes, decode)
	}

	encode = core.Base64URLEncode(bytes)
	decode, err = core.Base64URLDecode(encode)
	if err != nil {
		t.Fatal("Base64StdDecode error")
	}
	if !bytes2.Equal(decode, bytes) {
		t.Fatalf("Base64StdDecode error: want is %s, but acutal is %s", bytes, decode)
	}

	encode = core.Base64RawURLEncode(bytes)
	decode, err = core.Base64RawURLDecode(encode)
	if err != nil {
		t.Fatal("Base64RawURLDecode error")
	}
	if !bytes2.Equal(decode, bytes) {
		t.Fatalf("Base64RawURLDecode error: want is %s, but acutal is %s", bytes, decode)
	}

}
