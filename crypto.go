package gone

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io"

	"golang.org/x/crypto/ssh"

	"github.com/tjfoc/gmsm/sm2"
	x509sm "github.com/tjfoc/gmsm/x509"
)

// 最重要的一句话，公钥加密私钥解密，私钥签名公钥验证。

type RSABit int

const (
	RSA1024 RSABit = 1024
	RSA2048 RSABit = 2048
)

// GenerateBase64Key 生成base64编码的公私钥.
func GenerateBase64Key(secretLength SecretKeyLengthType,
	secretFormat SecretKeyFormatType) (pkStr, pbkStr string, err error) {
	var (
		privateKeyBytes []byte
		publicKeyBytes  []byte
		privateKey      *sm2.PrivateKey
		pkBytes         []byte
	)

	if secretLength == M2 {
		privateKey, err = sm2.GenerateKey(rand.Reader)
		if err != nil {
			return "", "", err
		}
		privateKeyBytes, err = x509sm.MarshalSm2UnecryptedPrivateKey(privateKey)
		if err != nil {
			return "", "", err
		}
		publicKeyBytes, err = x509sm.MarshalSm2PublicKey(&privateKey.PublicKey)
		if err != nil {
			return "", "", err
		}

		return base64.RawURLEncoding.EncodeToString(privateKeyBytes),
			base64.RawURLEncoding.EncodeToString(publicKeyBytes), nil
	}
	var priKey *rsa.PrivateKey
	if secretLength == RSA {
		priKey, err = rsa.GenerateKey(rand.Reader, 2048)
	}

	if secretFormat == PKCS1 {
		// 生成公匙
		pkStr = base64.RawURLEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(priKey))
		pbkStr = base64.RawURLEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&(priKey.PublicKey)))
		return pkStr, pbkStr, err
	}

	if secretFormat == PKCS8 {
		// 生成公匙
		pkBytes, err = x509.MarshalPKCS8PrivateKey(priKey)
		if err != nil {
			return "", "", err
		}
		pkStr = base64.RawURLEncoding.EncodeToString(pkBytes)
		pbkStr = base64.RawURLEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(&(priKey.PublicKey)))
		return pkStr, pbkStr, err
	}
	return "", "", err
}

// GenerateECDSAKey 生成PEM格式ECDSA公私钥，写入到io.Writer中.
func GenerateECDSAKey(privateWriter, publicWriter io.Writer, c elliptic.Curve) error {
	var (
		pk       *ecdsa.PrivateKey
		pkBytes  []byte
		pubBytes []byte
		err      error
	)
	pk, err = ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		panic(err)
	}
	pkBytes, err = x509.MarshalECPrivateKey(pk)
	if err != nil {
		return err
	}
	err = pem.Encode(privateWriter, &pem.Block{
		Type:  "ECD PRIVATE KEY",
		Bytes: pkBytes,
	})
	if err != nil {
		return err
	}
	pubBytes, err = x509.MarshalPKIXPublicKey(&pk.PublicKey)
	if err != nil {
		return err
	}
	return pem.Encode(publicWriter, &pem.Block{
		Type:  "ECD PUBLIC KEY",
		Bytes: pubBytes,
	})
}

// GenerateECDSAKeyToMemory 生成PEM格式ECDSA公私钥，返回字节格式.
func GenerateECDSAKeyToMemory(c elliptic.Curve) (privateBytes []byte, publicBytes []byte, err error) {
	privateBuffer := bytes.Buffer{}
	publicBuffer := bytes.Buffer{}
	err = GenerateECDSAKey(&privateBuffer, &publicBuffer, c)
	if err != nil {
		return privateBytes, publicBytes, err
	}
	privateBytes = privateBuffer.Bytes()
	publicBytes = publicBuffer.Bytes()
	return privateBytes, publicBytes, err
}

// GenerateSSHKey 生成ssh密钥队.
func GenerateSSHKey(bits RSABit) (pkBytes []byte, pbkBytes []byte, err error) {
	var priKey *rsa.PrivateKey
	priKey, err = rsa.GenerateKey(rand.Reader, int(bits))
	if err != nil {
		return nil, nil, err
	}
	pkBytes, err = x509.MarshalPKCS8PrivateKey(priKey)
	if err != nil {
		return nil, nil, err
	}
	// golang.org/x/crypto/ssh/keys.go:1129
	// PKCS8 使用pem编码使用PRIVATE KEY作为类型
	pkBytes = pem.EncodeToMemory(&pem.Block{
		Bytes: pkBytes,
		Type:  "PRIVATE KEY",
	})
	publicKey, err := ssh.NewPublicKey(&(priKey.PublicKey))
	if err != nil {
		return nil, nil, err
	}
	pbkBytes = ssh.MarshalAuthorizedKey(publicKey)
	return pkBytes, pbkBytes, err
}

var pemStart = []byte("\n-----BEGIN ")

// DecodePemHexBase64 解析pem或者hex或者base64编码成der编码.
func DecodePemHexBase64(data []byte) ([]byte, error) {
	if bytes.HasPrefix(data, pemStart[1:]) ||
		bytes.HasPrefix(data, pemStart) {
		block, _ := pem.Decode(data)
		if block == nil {
			return nil, errors.New("unable to decode publicKey to request")
		}
		return block.Bytes, nil
	}
	var err error
	decoded := make([]byte, hex.DecodedLen(len(data)))
	b, err := hex.Decode(decoded, data)
	// if parse ok return derBytes
	if err == nil {
		return decoded[:b], nil
	}
	decoded, err = Base64StdDecode(data)
	if err == nil {
		return decoded, nil
	}
	decoded, err = Base64RawStdDecode(data)
	if err == nil {
		return decoded, nil
	}
	decoded, err = Base64URLDecode(data)
	if err == nil {
		return decoded, nil
	}
	return Base64RawURLDecode(data)
}
