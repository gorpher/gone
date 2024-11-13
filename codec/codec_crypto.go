package codec

import (
	"crypto/cipher"
	"errors"
	crypto2 "github.com/gorpher/gone/cryptoutil"
)

type CryptoCodec interface {
	codec
	Encode(key, plaintext []byte) ([]byte, error)
	Decode(key, ciphertext []byte) ([]byte, error)
}

type cookieCodec struct {
	hashKey   []byte
	blockKey  []byte
	block     cipher.Block
	blockMode crypto2.BlockStreamMode
	maxLength int
	maxAge    int64
	minAge    int64
	err       error
}

var (
	errHashKeyNotSet     = errors.New("hash key is not set")
	errBlockKeyNotSet    = errors.New("block key is not set")
	errPlaintextInvalid  = errors.New("the plaintext is not valid")
	errPlaintextTooLong  = errors.New("the plaintext is too long")
	errCiphertextTooLong = errors.New("the ciphertext is too long")
	errCiphertextInvalid = errors.New("the ciphertext is not valid")
	errTimestampInvalid  = errors.New("invalid timestamp")
	errTimestampExpired  = errors.New("expired timestamp")
)
