package cryptoutil

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"strings"

	"github.com/tjfoc/gmsm/sm2"
	x509sm "github.com/tjfoc/gmsm/x509"
)

// SignBySM2Bytes 使用sm2私钥签名字符串，返回base64编码的license.
func SignBySM2Bytes(privateKey, licenseBytes []byte) (license string, err error) {
	var key *sm2.PrivateKey
	privateKey, err = base64.RawURLEncoding.DecodeString(string(privateKey))
	if err != nil {
		return "", err
	}
	key, err = x509sm.ParsePKCS8PrivateKey(privateKey, nil)
	if err != nil {
		return "", err
	}
	return SignBySM2(key, licenseBytes)
}

// SignBySM2  使用sm2私钥对象指针签名字符串，返回base64编码的license.
func SignBySM2(privateKey *sm2.PrivateKey, licenseBytes []byte) (license string, err error) {
	var (
		signBytes        []byte
		licenseBase64Str string
	)

	// 将授权信息json编码成base64字符串
	licenseBase64Str = base64.RawURLEncoding.EncodeToString(licenseBytes)

	// 用私钥对授权信息的base64字符串进行签名
	signBytes, err = privateKey.Sign(rand.Reader, []byte(licenseBase64Str), nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 将签名数据编码成base64字符串
	signBase64Str := base64.RawURLEncoding.EncodeToString(signBytes)
	// 拼接完整授权字符串
	license = licenseBase64Str + "." + signBase64Str
	return
}

// VerifyBySM2 使用sm2公钥验证签名的license.
func VerifyBySM2(publicKeyBase64, licenseCode string) (license string, valid bool, err error) {
	var (
		publicKeyBytes, signBytes []byte
		publicKey                 *sm2.PublicKey
		licenseInfo               []byte
	)

	// 解析公钥
	publicKeyBytes, err = base64.RawURLEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return
	}
	publicKey, err = x509sm.ParseSm2PublicKey(publicKeyBytes)
	if err != nil {
		return
	}

	// 从授权码中拆解出授权信息
	arr := strings.Split(licenseCode, ".")
	if len(arr) != 2 {
		err = errors.New("valid licenseCode")
		return
	}

	// 验证签名(被签名内容，签名)
	signBytes, err = base64.RawURLEncoding.DecodeString(arr[1])
	if err != nil {
		return
	}
	valid = publicKey.Verify([]byte(arr[0]), signBytes)
	if !valid {
		return
	}

	// 解析授权信息
	licenseInfo, err = base64.RawURLEncoding.DecodeString(arr[0])
	if err != nil {
		return "", false, err
	}
	return string(licenseInfo), true, nil
}

// SM2PublicEncrypt sm2公钥加密，参数publicKeyStr必须是hex、base64或者是pem编码.
func SM2PublicEncrypt(publicKeyBytes, textBytes []byte) ([]byte, error) {
	var (
		err       error
		publicKey *sm2.PublicKey
	)
	publicKeyBytes, err = DecodePemHexBase64(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	publicKey, err = ParseSM2PublicKey(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	return sm2.Encrypt(publicKey, textBytes, rand.Reader, sm2.C1C3C2)
}

// ParseSM2PublicKey 解析公钥，derBytes可以使用DecodePemHexBase64函数获取.
func ParseSM2PublicKey(derBytes []byte) (publicKey *sm2.PublicKey, err error) {
	var (
		pub interface{}
		ok  bool
	)
	publicKey, err = x509sm.ParseSm2PublicKey(derBytes)
	if err == nil {
		return publicKey, nil
	}
	err = nil //nolint

	pub, err = x509sm.ParsePKIXPublicKey(derBytes)
	if err != nil {
		return nil, errors.New("parse sm2 public key error:" + err.Error())
	}
	publicKey, ok = pub.(*sm2.PublicKey)
	if !ok {
		return nil, errors.New("failed to assert SM2 PublicKey type")
	}
	return publicKey, nil
}

// DecryptBySM2 使用SM2私钥解密.
func DecryptBySM2(privateKey *sm2.PrivateKey, ciphertext []byte) ([]byte, error) {
	return sm2.Decrypt(privateKey, ciphertext, sm2.C1C3C2)
}

func DecryptBySM2Bytes(privateKey []byte, ciphertext []byte) (text []byte, err error) {
	var key *sm2.PrivateKey
	privateKey, err = base64.RawURLEncoding.DecodeString(string(privateKey))
	if err != nil {
		return text, err
	}
	key, err = x509sm.ParsePKCS8PrivateKey(privateKey, nil)
	if err != nil {
		return text, err
	}
	return DecryptBySM2(key, ciphertext)
}
