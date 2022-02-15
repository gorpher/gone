package gone

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"math/big"
)

var (

	// ErrECDSANilPubKey is the error for trying to verify a JWT with a nil public key.
	ErrECDSANilPubKey = errors.New("jwt: ECDSA public key is nil")
	// ErrECDSAVerification is the error for an invalid ECDSA signature.
	ErrECDSAVerification = errors.New("jwt: ECDSA verification failed")

	// ErrRSANilPubKey is the error for trying to verify a JWT with a nil public key.
	ErrRSANilPubKey = errors.New("jwt: RSA public key is nil")
	// ErrRSAVerification is the error for an invalid RSA signature.
	ErrRSAVerification = errors.New("jwt: RSA verification failed")

	// ErrHmacVerification is the error for an invalid RSA signature.
	ErrHmacVerification = errors.New("jwt: Hmac verification failed")
)

type JwtPayload struct {
	JTI       string        `json:"jti"`
	IAT       int           `json:"iat"`
	EXP       int           `json:"exp"`
	NBF       int           `json:"nbf"`
	ISS       string        `json:"iss"`
	Sub       string        `json:"sub"`
	Aud       string        `json:"aud"`
	Roles     []interface{} `json:"roles"`
	AuthType  string        `json:"authType"`
	UserId    string        `json:"userId"`
	SubjectId string        `json:"subjectId"`
	Email     string        `json:"email"`
}

type JWK struct {
	KTY       string `json:"kty"`
	KID       string `json:"kid"`
	Use       string `json:"use"`
	Algorithm string `json:"alg"`
	N         string `json:"n"`
	E         string `json:"e"`
}

func RsaPublicKeyToJWK(pub *rsa.PublicKey) ([]byte, error) {
	return json.Marshal(JWK{
		E:         base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
		N:         base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
		KTY:       "RSA",
		Algorithm: "RS256",
		KID:       "0",
	})
}

func ParseRsaPublicKeyByJWK(jsonBytes []byte) (publicKey *rsa.PublicKey, err error) {
	var pks JWK
	err = json.Unmarshal(jsonBytes, &pks)
	if err != nil {
		return nil, err
	}
	if pks.KTY != "RSA" {
		return nil, errors.New("Unknown key type algorithm: '" + pks.KTY + "'")
	}
	n, err := Base64RawURLDecode([]byte(pks.N))
	if err != nil {
		return nil, err
	}
	e, err := Base64RawURLDecode([]byte(pks.E))
	if err != nil {
		return nil, err
	}
	bigN := new(big.Int)
	bigE := new(big.Int)
	bigN.SetBytes(n)
	bigE.SetBytes(e)
	intE := bigE.Int64()
	return &rsa.PublicKey{
		N: bigN,
		E: int(intE),
	}, nil
}

func jwtCryptoByHash(alg string) (func() hash.Hash, error) {
	switch alg {
	case "RS256":
		return crypto.SHA256.New, nil
	case "HS384":
		return crypto.SHA384.New, nil
	case "RS512":
		return crypto.SHA512.New, nil
	case "PS256":
		return crypto.SHA512.New, nil
	case "PS384":
		return crypto.SHA512.New, nil
	case "PS512":
		return crypto.SHA512.New, nil
	case "ES256":
		return crypto.SHA512.New, nil
	case "ES384":
		return crypto.SHA512.New, nil
	case "ES512":
		return crypto.SHA512.New, nil
	default:
		return nil, fmt.Errorf("the %s algorithm is not supported", alg)
	}
}

func VerifyJwtSignByRsa(ciphertext []byte, publicKey *rsa.PublicKey) (json.RawMessage, error) {
	sep1 := bytes.IndexByte(ciphertext, '.')
	if sep1 < 0 {
		return nil, ErrMalformed
	}
	cbytes := ciphertext[sep1+1:]
	sep2 := bytes.IndexByte(cbytes, '.')
	if sep2 < 0 {
		return nil, ErrMalformed
	}
	sign := cbytes[sep2+1:]
	sep3 := bytes.IndexByte(cbytes, '.')
	if sep3 < 0 {
		return nil, ErrMalformed
	}
	hb, err := base64.RawURLEncoding.DecodeString(string(ciphertext[:sep1]))
	if err != nil {
		return nil, err
	}
	var header Header
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(hb, &header)
	if err != nil {
		return nil, err
	}
	hs, err := jwtCryptoByHash(header.Algorithm)
	if err != nil {
		return nil, err
	}
	signBytes, err := base64.RawURLEncoding.DecodeString(string(sign))
	if err != nil {
		return nil, err
	}
	h := hs()
	p := ciphertext[:sep1+1+sep2]
	h.Write(p)
	if header.Algorithm == "RS256" {
		err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, h.Sum(nil), signBytes)
		if err != nil {
			return nil, err
		}
	}
	if header.Algorithm == "RS384" {
		err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA384, h.Sum(nil), signBytes)
		if err != nil {
			return nil, err
		}
	}
	if header.Algorithm == "RS512" {
		err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA512, h.Sum(nil), signBytes)
		if err != nil {
			return nil, err
		}
	}
	if header.Algorithm == "PS256" {
		err := rsa.VerifyPSS(publicKey, crypto.SHA256, h.Sum(nil), signBytes, &rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthAuto,
			Hash:       crypto.SHA256,
		})
		if err != nil {
			return nil, err
		}
	}
	if header.Algorithm == "PS384" {
		err := rsa.VerifyPSS(publicKey, crypto.SHA384, h.Sum(nil), signBytes, &rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthAuto,
			Hash:       crypto.SHA384,
		})
		if err != nil {
			return nil, err
		}
	}
	if header.Algorithm == "PS512" {
		err := rsa.VerifyPSS(publicKey, crypto.SHA512, h.Sum(nil), signBytes, &rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthAuto,
			Hash:       crypto.SHA512,
		})
		if err != nil {
			return nil, err
		}
	}
	return base64.RawURLEncoding.DecodeString(string(cbytes[:sep2]))
}

func VerifyJwtSignByEcdsa(ciphertext []byte, publicKey *ecdsa.PublicKey) (json.RawMessage, error) {
	var body json.RawMessage
	sep1 := bytes.IndexByte(ciphertext, '.')
	if sep1 < 0 {
		return nil, ErrMalformed
	}
	cbytes := ciphertext[sep1+1:]
	sep2 := bytes.IndexByte(cbytes, '.')
	if sep2 < 0 {
		return nil, ErrMalformed
	}
	sign := cbytes[sep2+1:]
	sep3 := bytes.IndexByte(cbytes, '.')
	if sep3 < 0 {
		return nil, ErrMalformed
	}
	hb, err := base64.RawURLEncoding.DecodeString(string(ciphertext[:sep1]))
	if err != nil {
		return nil, err
	}
	var header Header
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(hb, &header)
	if err != nil {
		return nil, err
	}
	hs, err := jwtCryptoByHash(header.Algorithm)
	if err != nil {
		return nil, err
	}
	signBytes, err := base64.RawURLEncoding.DecodeString(string(sign))
	if err != nil {
		return nil, err
	}
	h := hs()
	p := ciphertext[:sep1+1+sep2]
	h.Write(p)
	byteSize := byteSize(publicKey.Params().BitSize)
	if len(signBytes) != byteSize*2 {
		return nil, ErrECDSAVerification
	}
	r := big.NewInt(0).SetBytes(signBytes[:byteSize])
	s := big.NewInt(0).SetBytes(signBytes[byteSize:])
	if !ecdsa.Verify(publicKey, h.Sum(nil), r, s) {
		return nil, ErrECDSAVerification
	}
	body, err = base64.RawURLEncoding.DecodeString(string(cbytes[:sep2]))
	if err != nil {
		return nil, err
	}
	return body, nil
}

func VerifyJwtSignByHmacHash(ciphertext, key []byte) (json.RawMessage, error) {
	sep1 := bytes.IndexByte(ciphertext, '.')
	if sep1 < 0 {
		return nil, ErrMalformed
	}
	cbytes := ciphertext[sep1+1:]
	sep2 := bytes.IndexByte(cbytes, '.')
	if sep2 < 0 {
		return nil, ErrMalformed
	}

	sig := cbytes[sep2+1:]
	sep3 := bytes.IndexByte(cbytes, '.')
	if sep3 < 0 {
		return nil, ErrMalformed
	}
	hb, err := base64.RawURLEncoding.DecodeString(string(ciphertext[:sep1]))
	if err != nil {
		return nil, err
	}

	var header Header
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(hb, &header)
	if err != nil {
		return nil, err
	}

	hs, err := jwtHS(header.Algorithm)
	if err != nil {
		return nil, err
	}
	h := hs(key)
	p := ciphertext[:sep1+1+sep2]
	h.Write(p)

	sha256 := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	if sha256 != string(sig) {
		return nil, errors.New("invalid jwt token")
	}
	return base64.RawURLEncoding.DecodeString(string(cbytes[:sep2]))
}

// VerifyJwtSign verify jwt sign text , key is hmac key or rsa/ecdsa public key
func VerifyJwtSign(ciphertext, key []byte) (json.RawMessage, error) {
	sep1 := bytes.IndexByte(ciphertext, '.')
	if sep1 < 0 {
		return nil, ErrMalformed
	}
	cbytes := ciphertext[sep1+1:]
	sep2 := bytes.IndexByte(cbytes, '.')
	if sep2 < 0 {
		return nil, ErrMalformed
	}
	sign := cbytes[sep2+1:]
	sep3 := bytes.IndexByte(cbytes, '.')
	if sep3 < 0 {
		return nil, ErrMalformed
	}
	hb, err := base64.RawURLEncoding.DecodeString(string(ciphertext[:sep1]))
	if err != nil {
		return nil, err
	}
	var header Header
	err = json.Unmarshal(hb, &header)
	if err != nil {
		return nil, err
	}
	signBytes, err := base64.RawURLEncoding.DecodeString(string(sign))
	if err != nil {
		return nil, err
	}
	if len(header.Algorithm) < 5 {
		return nil, fmt.Errorf("the %s algorithm is not supported", header.Algorithm)
	}
	alg := header.Algorithm[:2]
	size := header.Algorithm[2:]
	if !(alg == "HS" || alg == "RS" || alg == "ES") || !(size == "256" || size == "284" || size == "512") {
		return nil, fmt.Errorf("the %s algorithm is not supported", header.Algorithm)
	}

	var h hash.Hash
	var cryptoHash crypto.Hash
	if size == "256" {
		if alg == "HS" {
			h = hmac.New(crypto.SHA256.New, key)
		} else {
			h = crypto.SHA256.New()
		}
		cryptoHash = crypto.SHA256
	}
	if size == "384" {
		if alg == "HS" {
			h = hmac.New(crypto.SHA384.New, key)
		} else {
			h = crypto.SHA384.New()
		}
		cryptoHash = crypto.SHA384
	}
	if size == "512" {
		if alg == "HS" {
			h = hmac.New(crypto.SHA512.New, key)
		} else {
			h = crypto.SHA512.New()
		}
		cryptoHash = crypto.SHA512
	}

	p := ciphertext[:sep1+1+sep2]
	h.Write(p)

	if alg == "RS" || alg == "PS" {
		var (
			rsaPublicKey *rsa.PublicKey
			pbk          interface{}
		)
		rsaPublicKey, err = x509.ParsePKCS1PublicKey(key)
		if err != nil {
			pbk, err = x509.ParsePKIXPublicKey(key)
			if err != nil {
				return nil, err
			}
			var ok bool
			rsaPublicKey, ok = pbk.(*rsa.PublicKey)
			if !ok {
				return nil, ErrRSANilPubKey
			}
		}
		if alg == "RS" {
			err = rsa.VerifyPKCS1v15(rsaPublicKey, cryptoHash, h.Sum(nil), signBytes)
			if err != nil {
				return nil, err
			}
		}

		if alg == "PS" {
			err = rsa.VerifyPSS(rsaPublicKey, cryptoHash, h.Sum(nil), signBytes, &rsa.PSSOptions{
				SaltLength: rsa.PSSSaltLengthAuto,
				Hash:       cryptoHash,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	if alg == "ES" {
		var (
			ecdsaPublicKey *ecdsa.PublicKey
			pbk            interface{}
		)
		pbk, err = x509.ParsePKIXPublicKey(key)
		if err != nil {
			return nil, err
		}
		var ok bool
		ecdsaPublicKey, ok = pbk.(*ecdsa.PublicKey)
		if !ok {
			return nil, ErrECDSANilPubKey
		}
		byteSize := byteSize(ecdsaPublicKey.Params().BitSize)
		if len(signBytes) != byteSize*2 {
			return nil, ErrECDSAVerification
		}
		r := big.NewInt(0).SetBytes(signBytes[:byteSize])
		s := big.NewInt(0).SetBytes(signBytes[byteSize:])
		if !ecdsa.Verify(ecdsaPublicKey, h.Sum(nil), r, s) {
			return nil, ErrECDSAVerification
		}
	}

	if alg == "HS" {
		sha256 := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
		if sha256 != string(sign) {
			return nil, ErrHmacVerification
		}
	}
	return base64.RawURLEncoding.DecodeString(string(cbytes[:sep2]))
}

func byteSize(bitSize int) int {
	byteSize := bitSize / 8
	if bitSize%8 > 0 {
		return byteSize + 1
	}
	return byteSize
}
