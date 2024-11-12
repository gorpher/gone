package crypto

import (
	"fmt"
	"github.com/gorpher/gone/core"
	"testing"
)

func TestSignBySM2Bytes(t *testing.T) {
	pkStr, pbkStr, err := GenerateBase64Key(core.M2, core.PKCS1)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(pkStr, "\n", pbkStr)

	signStr, err := SignBySM2Bytes([]byte(pkStr), []byte(licenseStr))
	if err != nil {
		t.Error(err)
		return
	}
	rsaStr, valid, err := VerifyBySM2(pbkStr, signStr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(rsaStr, valid)

	//=============================================================================
	pkStr, pbkStr, err = GenerateBase64Key(core.M2, core.PKCS8)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(pkStr, "\n", pbkStr)
	signStr, err = SignBySM2Bytes([]byte(pkStr), []byte(licenseStr))
	if err != nil {
		t.Error(err)
		return
	}
	rsaStr, valid, err = VerifyBySM2(pbkStr, signStr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(rsaStr, valid)
}

func TestSM2PublicEncrypt(t *testing.T) {
	pkStr, pbkStr, err := GenerateBase64Key(core.M2, core.PKCS1)
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
