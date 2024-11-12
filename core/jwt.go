package core

import (
	"crypto"
	"crypto/hmac"
	"fmt"
	"hash"
)

func JwtHS(alg string) (func(key []byte) hash.Hash, error) {
	switch alg {
	case "HS256":
		return func(key []byte) hash.Hash {
			return hmac.New(crypto.SHA256.New, key)
		}, nil
	case "HS384":
		return func(key []byte) hash.Hash {
			return hmac.New(crypto.SHA384.New, key)
		}, nil
	case "HS512":
		return func(key []byte) hash.Hash {
			return hmac.New(crypto.SHA512.New, key)
		}, nil
	default:
		return nil, fmt.Errorf("the %s algorithm is not supported", alg)
	}
}
