package gone

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
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

type jwtCodec struct {
	err       error
	hashFunc  func([]byte) hash.Hash
	algorithm string
	sz        Serializer
}

type Header struct {
	Algorithm   string `json:"alg,omitempty"`
	ContentType string `json:"cty,omitempty"`
	KeyID       string `json:"kid,omitempty"`
	Type        string `json:"typ,omitempty"`
}

// Payload is a JWT payload according to the RFC 7519.
type Payload struct {
	Issuer         string   `json:"iss,omitempty"`
	Subject        string   `json:"sub,omitempty"`
	Audience       Audience `json:"aud,omitempty"`
	ExpirationTime *Time    `json:"exp,omitempty"`
	NotBefore      *Time    `json:"nbf,omitempty"`
	IssuedAt       *Time    `json:"iat,omitempty"`
	JWTID          string   `json:"jti,omitempty"`
}

// Audience is a special claim that may either be
// a single string or an array of strings, as per the RFC 7519.
type Audience []string

// MarshalJSON implements a marshaling function for "aud" claim.
func (a Audience) MarshalJSON() ([]byte, error) {
	switch len(a) {
	case 0:
		return json.Marshal("") // nil or empty slice returns an empty string
	case 1:
		return json.Marshal(a[0])
	default:
		return json.Marshal([]string(a))
	}
}

// UnmarshalJSON implements an unmarshaling function for "aud" claim.
func (a *Audience) UnmarshalJSON(b []byte) error {
	var (
		v   interface{}
		err error
	)
	if err = json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch vv := v.(type) {
	case string:
		aud := make(Audience, 1)
		aud[0] = vv
		*a = aud
	case []interface{}:
		aud := make(Audience, len(vv))
		for i := range vv {
			aud[i] = vv[i].(string)
		}
		*a = aud
	}
	return nil
}

func NewJwtCodec(alg string) Codec {
	j := &jwtCodec{
		sz: JSONEncoder{},
	}

	hashFunc, err := jwtHS(alg)
	j.hashFunc = hashFunc
	j.algorithm = alg
	j.err = err
	return j
}
func jwtHS(alg string) (func(key []byte) hash.Hash, error) {
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
		return nil, errors.New("the algorithm is not supported")
	}
}

// ErrNotJSONObject is the error for when a JWT payload is not a JSON object.
var ErrNotJSONObject = errors.New("jwt: payload is not a valid JSON object")

// ErrMalformed indicates a token doesn't have a valid format, as per the RFC 7519.
var ErrMalformed = errors.New("jwt: malformed token")

func isJSONObject(payload []byte) bool {
	payload = bytes.TrimSpace(payload)
	return payload[0] == '{' && payload[len(payload)-1] == '}'
}

func (j *jwtCodec) Encode(key, plaintext []byte) ([]byte, error) {
	if j.err != nil {
		return nil, j.err
	}
	hb, err := j.sz.Serialize(Header{
		Algorithm: j.algorithm,
		Type:      "JWT",
	})
	if err != nil {
		return nil, err
	}
	if !isJSONObject(plaintext) {
		return nil, ErrNotJSONObject
	}
	h := j.hashFunc(key)
	enc := base64.RawURLEncoding
	h64len := enc.EncodedLen(len(hb))
	p64len := enc.EncodedLen(len(plaintext))
	sig64len := enc.EncodedLen(h.Size())
	token := make([]byte, h64len+1+p64len+1+sig64len)
	enc.Encode(token, hb)
	token[h64len] = '.'
	enc.Encode(token[h64len+1:], plaintext)
	token[h64len+1+p64len] = '.'
	h.Write(token[:h64len+1+p64len])
	sig := h.Sum(nil)
	enc.Encode(token[h64len+1+p64len+1:], sig)
	return token, nil
}

func (j *jwtCodec) Decode(key, ciphertext []byte) ([]byte, error) {
	if j.err != nil {
		return nil, j.err
	}
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
	hp, err := base64.RawURLEncoding.DecodeString(string(ciphertext[:sep1]))
	if err != nil {
		return nil, err
	}

	var header Header
	err = j.sz.Deserialize(hp, &header)
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
	return cbytes[sep2+1:], nil
}
