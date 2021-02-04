package gone

import (
	"testing"
)

func TestAesEncryptCBC(t *testing.T) {
	key := []byte(RandString(16))
	iv := []byte(RandString(16))
	txt := []byte(RandString(1024))
	encrypt, err := AesEncryptCBC(txt, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	taed, err := AesDecryptCBC(encrypt, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	if len(txt) != len(taed) {
		t.Fatalf("原文和解密后的原文不匹配 %s!=%s", txt, taed)
	}
	for i := range txt {
		if txt[i] != taed[i] {
			t.Fatalf("原文和解密后的原文不匹配 %s!=%s", txt, taed)
		}
	}
}
