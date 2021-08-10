package gone

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Codec interface {
	Encode(key, plaintext []byte) ([]byte, error)
	Decode(key, ciphertext []byte) ([]byte, error)
}

type cookieCodec struct {
	hashKey   []byte
	blockKey  []byte
	block     cipher.Block
	blockMode BlockStreamMode
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

func NewCookieCodec(hashKey, blockKey []byte) Codec {
	s := &cookieCodec{
		hashKey:   hashKey,
		blockKey:  blockKey,
		maxAge:    86400 * 30,
		maxLength: 4096,
		blockMode: CTR,
	}
	if hashKey == nil {
		s.err = errHashKeyNotSet
	}
	if blockKey == nil {
		s.err = errBlockKeyNotSet
	}
	s.block, s.err = aes.NewCipher(s.blockKey)
	return s
}

func (s *cookieCodec) Encode(key, plaintext []byte) ([]byte, error) {
	if s.err != nil {
		return nil, s.err
	}
	if len(plaintext) == 0 {
		return nil, errPlaintextInvalid
	}
	var err error
	var b = make([]byte, len(plaintext))
	copy(b, plaintext)
	// 1. 加密.
	b, err = BlockEncrypt(s.block, s.blockMode, b)
	if err != nil {
		return nil, err
	}
	b = Base64URLEncode(b)
	// 2. 根据 "key|date|plaintext" 格式生成hash.
	b = []byte(fmt.Sprintf("%s|%d|%s|", key, s.timestamp(), b))
	mac := HMacSha256(s.hashKey, b[:len(b)-1])
	// 删除key，尾部追加hash
	b = append(b, mac...)[len(key)+1:]
	// 3. 使用base64编码.
	b = Base64URLEncode(b)
	// 4. 检查长度
	if s.maxLength != 0 && len(b) > s.maxLength {
		return nil, errPlaintextTooLong
	}
	return b, nil
}

func (s *cookieCodec) timestamp() int64 {
	return time.Now().UTC().Unix()
}

func (s *cookieCodec) Decode(key, ciphertext []byte) ([]byte, error) {
	if s.err != nil {
		return nil, s.err
	}
	// 1. 检查长度.
	if s.maxLength != 0 && len(ciphertext) > s.maxLength {
		return nil, errCiphertextTooLong
	}
	// 2. 解码base64.
	b, err := Base64URLDecode(ciphertext)
	if err != nil {
		return nil, err
	}
	// 3. 验证hash值. 格式： "date|ciphertext|mac".
	parts := bytes.SplitN(b, []byte("|"), 3)
	if len(parts) != 3 {
		return nil, errCiphertextInvalid
	}

	b = append([]byte(string(key)+"|"), b[:len(b)-len(parts[2])-1]...)
	if string(parts[2]) != HMacSha256(s.hashKey, b) {
		return nil, errCiphertextInvalid
	}
	// 4. 验证日期范围.
	var t1 int64
	if t1, err = strconv.ParseInt(string(parts[0]), 10, 64); err != nil {
		return nil, errTimestampInvalid
	}
	t2 := s.timestamp()
	if s.minAge != 0 && t1-t2 < s.minAge {
		return nil, errTimestampExpired
	}
	if s.maxAge != 0 && t1-t2 > s.maxAge {
		return nil, errTimestampExpired
	}
	// 5. 解密.
	b, err = Base64URLDecode(parts[1])
	if err != nil {
		return nil, err
	}
	return BlockDecrypt(s.block, s.blockMode, b)
}
