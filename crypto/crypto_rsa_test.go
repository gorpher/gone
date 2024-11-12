package crypto

import (
	bytes2 "bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gorpher/gone/core"
	"os"
	"testing"
)

func TestGenerateRSAKey(t *testing.T) {
	err := GenerateRSAKey(os.Stdout, os.Stdout, RSA1024)
	if err != nil {
		t.Error(err)
		return
	}
	err = GenerateRSAKey(os.Stdout, os.Stdout, RSA2048)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestEncryptByRSABytes(t *testing.T) {
	pkBytes, pbkBytes, err := GenerateRSAKeyToMemory(2048)
	if err != nil {
		t.Error(err)
		return
	}
	var publicKey *rsa.PublicKey
	pbkBlock, _ := pem.Decode(pbkBytes)
	if pbkBlock == nil {
		t.Error(errors.New("public key error"))
		return
	}
	publicKey, err = x509.ParsePKCS1PublicKey(pbkBlock.Bytes)
	if err != nil {
		t.Error(err)
		return
	}
	var privateKey *rsa.PrivateKey
	pkBlock, _ := pem.Decode(pkBytes)
	if pkBlock == nil {
		t.Error(errors.New("private key error"))
		return
	}
	var pi interface{}
	var ok bool
	pi, err = x509.ParsePKCS1PrivateKey(pkBlock.Bytes)
	if err != nil {
		t.Error(err)
		return
	}
	privateKey, ok = pi.(*rsa.PrivateKey)
	if !ok || privateKey == nil {
		t.Error(errors.New("private key assert failed"))
	}
	content := []byte("hello world")
	bytes, err := EncryptByRSA(publicKey, content)
	if err != nil {
		t.Error(err)
		return
	}
	decryptBytes, err := DecryptByRSA(privateKey, bytes)
	if err != nil {
		t.Error(err)
		return
	}
	if !bytes2.Equal(decryptBytes, content) {
		t.Error("DecryptByRSA ERROR")
	}
}

func TestSignByRSABytes(t *testing.T) {
	var (
		pkStr   string
		pbkStr  string
		signStr string
		license string
		valid   bool
		err     error
	)
	//=============================================================================
	pkStr, pbkStr, err = GenerateBase64Key(core.RSA, core.PKCS1)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(pkStr, "\n", pbkStr)

	signStr, err = SignByRSABytes([]byte(pkStr), []byte(licenseStr))
	if err != nil {
		t.Error(err)
		return
	}
	license, valid, err = VerifyByRSA(pbkStr, signStr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(license, valid)

	//=============================================================================
	pkStr, pbkStr, err = GenerateBase64Key(core.RSA, core.PKCS8)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(pkStr, "\n", pbkStr)

	signStr, err = SignByRSABytes([]byte(pkStr), []byte(licenseStr))
	if err != nil {
		t.Error(err)
		return
	}
	license, valid, err = VerifyByRSA(pbkStr, signStr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(license, valid)
}

func TestParsePublicKey(t *testing.T) {
	priKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	// pkcs1 解析
	key := x509.MarshalPKCS1PublicKey(&(priKey.PublicKey))
	_, err = ParseRsaPublicKey(key)
	if err != nil {
		t.Fatal(err)
	}
	// pkix 解析
	pbkBytes, err := x509.MarshalPKIXPublicKey(&(priKey.PublicKey))
	if err != nil {
		t.Fatal(err)
	}

	_, err = x509.ParsePKIXPublicKey(pbkBytes)
	if err != nil {
		t.Fatal(err)
	}
	// 自己的组合解析
	pbkBytes, err = DecodePemHexBase64([]byte(PublicKeyPemStr))
	if err != nil {
		t.Fatal(err)
	}
	_, err = ParseRsaPublicKey(pbkBytes)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParsePrivateKey(t *testing.T) {
	priKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	// pkcs1 解析
	_, err = ParseRsaPrivateKey(x509.MarshalPKCS1PrivateKey(priKey))
	if err != nil {
		t.Fatal(err)
	}

	// pkcs8 解析
	pkBytes, err := x509.MarshalPKCS8PrivateKey(priKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = x509.ParsePKCS8PrivateKey(pkBytes)
	if err != nil {
		t.Fatal(err)
	}
	pkBytes, err = DecodePemHexBase64([]byte(PrivateKeyPemStr))
	if err != nil {
		t.Fatal(err)
	}
	_, err = ParseRsaPrivateKey(pkBytes)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRsaPublicEncrypt(t *testing.T) {
	// 每次加密的字节数，不能超过密钥的长度值减去11,而每次加密得到的密文长度，却恰恰是密钥的长度
	textBytes := []byte(testRandStr(100))
	encrypt, err := RsaPublicEncrypt([]byte(PublicKeyPemStr), textBytes)
	if err != nil {
		t.Fatal(err)
	}
	decryptBytes, err := RsaPrivateDecrypt([]byte(PrivateKeyPemStr), encrypt)
	if err != nil {
		t.Fatal(err)
	}
	if len(decryptBytes) != len(textBytes) {
		t.Error("error")
	}
	if string(decryptBytes) != string(textBytes) {
		t.Error("error")
	}
	encrypt, err = RsaPublicEncrypt([]byte(PublicKeyBase64Str), textBytes)
	if err != nil {
		t.Fatal(err)
	}
	decryptBytes, err = RsaPrivateDecrypt([]byte(PrivateKeyBase64Str), encrypt)
	if err != nil {
		t.Fatal(err)
	}
	if len(decryptBytes) != len(textBytes) {
		t.Error("error")
	}
	if string(decryptBytes) != string(textBytes) {
		t.Error("error")
	}
	encrypt, err = RsaPublicEncrypt([]byte(PublicKeyHexStr), textBytes)
	if err != nil {
		t.Fatal(err)
	}
	decryptBytes, err = RsaPrivateDecrypt([]byte(PrivateKeyHexStr), encrypt)
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
