package gone

import (
	"fmt"
	"testing"
)

func TestSM2PublicEncrypt(t *testing.T) {
	pkStr, pbkStr, err := GenerateBase64Key(M2, PKCS1)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(pkStr, "\n", pbkStr)
	textBytes := []byte(testRandStr(100))
	encrypt, err := SM2PublicEncrypt([]byte(pbkStr), textBytes)
	if err != nil {
		t.Fatal(err)
	}
	decryptBytes, err := DecryptBySM2Bytes([]byte(pkStr), encrypt)
	if err != nil {
		t.Fatal(err)
	}
	if len(decryptBytes) != len(textBytes) {
		t.Error("error")
	}
	if string(decryptBytes) != string(textBytes) {
		t.Error("error")
	}

}
