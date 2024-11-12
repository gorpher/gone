package crypto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"log"
	"strings"
)

// GenerateRSAKeyToMemory 生成PEM格式PKCS1的RSA公私钥，返回字节格式.
func GenerateRSAKeyToMemory(bits RSABit) (privateBytes []byte, publicBytes []byte, err error) {
	privateBuffer := bytes.Buffer{}
	publicBuffer := bytes.Buffer{}
	err = GenerateRSAKey(&privateBuffer, &publicBuffer, bits)
	if err != nil {
		return privateBytes, publicBytes, err
	}
	privateBytes = privateBuffer.Bytes()
	publicBytes = publicBuffer.Bytes()
	return privateBytes, publicBytes, err
}

// GenerateRSAKey 生成PEM格式PKCS1的RSA公私钥，写入到io.Writer中.
func GenerateRSAKey(privateWriter, publicWriter io.Writer, bits RSABit) error {
	var priKey *rsa.PrivateKey
	var pbkBytes []byte
	var err error
	priKey, err = rsa.GenerateKey(rand.Reader, int(bits))
	if err != nil {
		return err
	}
	pbkBytes = x509.MarshalPKCS1PublicKey(&(priKey.PublicKey))
	err = pem.Encode(privateWriter, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priKey),
	})
	if err != nil {
		return err
	}
	return pem.Encode(publicWriter, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pbkBytes,
	})
}

// SignByRSABytes 使用rsa私钥签名字符串，返回base64编码的license.
func SignByRSABytes(key, licenseBytes []byte) (string, error) {
	var (
		privateKey *rsa.PrivateKey
		pri2       interface{}
		err        error
		ok         bool
	)
	key, err = base64.RawURLEncoding.DecodeString(string(key))
	if err != nil {
		return "", err
	}
	privateKey, err = x509.ParsePKCS1PrivateKey(key)
	if err == nil {
		return SignByRSA(privateKey, licenseBytes)
	}
	pri2, err = x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		return "", err
	}
	if privateKey, ok = pri2.(*rsa.PrivateKey); !ok {
		return "", errors.New("x509: failed to parse private key")
	}
	return SignByRSA(privateKey, licenseBytes)
}

// SignByRSA 使用rsa私钥对象指针签名字符串，返回base64编码的license.
func SignByRSA(key *rsa.PrivateKey, licenseBytes []byte) (license string, err error) {
	var (
		signBytes        []byte
		licenseBase64Str string
	)
	// 将授权信息json编码成base64字符串
	licenseBase64Str = base64.RawURLEncoding.EncodeToString(licenseBytes)
	hash := sha256.New()
	hash.Write([]byte(licenseBase64Str)) //nolint
	signBytes, err = key.Sign(rand.Reader, hash.Sum(nil), crypto.SHA256)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// 将签名数据编码成base64字符串
	signBase64Str := base64.RawURLEncoding.EncodeToString(signBytes)

	// 拼接完整授权字符串
	license = licenseBase64Str + "." + signBase64Str
	return license, nil
}

// VerifyByRSA 使用rsa公钥验证签名的license.
func VerifyByRSA(publicKeyBase64, licenseCode string) (license string, valid bool, err error) {
	var (
		publicKeyBytes, signBytes []byte
		publicKey                 *rsa.PublicKey
		licenseInfo               []byte
	)

	// 解析公钥
	publicKeyBytes, err = base64.RawURLEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return
	}
	publicKey, err = x509.ParsePKCS1PublicKey(publicKeyBytes)
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

	hash := sha256.New()
	hash.Write([]byte(arr[0])) //nolint

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash.Sum(nil), signBytes)
	if err != nil {
		log.Println(err)
		return "", false, err
	}

	// 解析授权信息
	licenseInfo, err = base64.RawURLEncoding.DecodeString(arr[0])
	if err != nil {
		return "", false, err
	}
	return string(licenseInfo), true, nil
}

// EncryptByRSABytes 使用RSA公钥加密.
func EncryptByRSABytes(publicKey, content []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pi, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	var pb *rsa.PublicKey
	var ok bool
	pb, ok = pi.(*rsa.PublicKey)
	if !ok || pb == nil {
		return nil, errors.New("public key assert failed")
	}
	return EncryptByRSA(pb, content)
}

// DecryptByRSABytes 使用RSA私钥解密.
func DecryptByRSABytes(privateKey []byte, ciphertext []byte) ([]byte, error) {
	var pk *rsa.PrivateKey
	var ok bool
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	pk, errv := x509.ParsePKCS1PrivateKey(block.Bytes)
	if errv != nil {
		pi, errv := x509.ParsePKCS8PrivateKey(block.Bytes)
		if errv != nil {
			return nil, errv
		}
		pk, ok = pi.(*rsa.PrivateKey)
		if !ok || pk == nil {
			return nil, errors.New("private key assert failed")
		}
	}
	return rsa.DecryptPKCS1v15(rand.Reader, pk, ciphertext)
}

// EncryptByRSA 使用RSA公钥加密.
func EncryptByRSA(publicKey *rsa.PublicKey, content []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, content)
}

// DecryptByRSA 使用RSA私钥解密.
func DecryptByRSA(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}

// RsaPublicEncrypt Rsa公钥加密，参数publicKeyStr必须是hex、base64或者是pem编码.
func RsaPublicEncrypt(publicKeyBytes, textBytes []byte) ([]byte, error) {
	var (
		err       error
		publicKey *rsa.PublicKey
	)
	publicKeyBytes, err = DecodePemHexBase64(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	publicKey, err = ParseRsaPublicKey(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, textBytes)
}

// ParseRsaPublicKey 解析公钥，derBytes可以使用DecodePemHexBase64函数获取.
func ParseRsaPublicKey(derBytes []byte) (publicKey *rsa.PublicKey, err error) {
	var (
		pub interface{}
		ok  bool
	)
	publicKey, err = x509.ParsePKCS1PublicKey(derBytes)
	if err == nil {
		return publicKey, nil
	}
	err = nil //nolint

	pub, err = x509.ParsePKIXPublicKey(derBytes)
	if err != nil {
		return nil, errors.New("parse rsa public key error: " + err.Error())
	}
	publicKey, ok = pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to assert RSA PublicKey type")
	}
	return publicKey, nil
}

// RsaPrivateDecrypt 解析rsa私钥，参数privateKeyStr必须是hex、base64或者是pem编码.
func RsaPrivateDecrypt(privateKeyBytes, cipherBytes []byte) (textBytes []byte, err error) {
	var privateKey *rsa.PrivateKey
	derBytes, err := DecodePemHexBase64(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	privateKey, err = ParseRsaPrivateKey(derBytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherBytes)
}

// ParseRsaPrivateKey 解析私钥，derBytes可以使用DecodePemHexBase64函数获取.
func ParseRsaPrivateKey(derBytes []byte) (privateKey *rsa.PrivateKey, err error) {
	var (
		pk interface{}
		ok bool
	)

	privateKey, err = x509.ParsePKCS1PrivateKey(derBytes)
	// if parse ok return private key
	if err == nil {
		return privateKey, nil
	}
	err = nil // nolint
	pk, err = x509.ParsePKCS8PrivateKey(derBytes)
	if err != nil {
		return nil, errors.New("parse rsa public key error:" + err.Error())
	}
	privateKey, ok = pk.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("failed to assert RSA PrivateKey type")
	}
	return privateKey, nil
}
