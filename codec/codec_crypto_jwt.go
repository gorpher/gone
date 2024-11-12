package codec

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gorpher/gone/core"
	"hash"
)

type jwtCodec struct {
	err       error
	hashFunc  func([]byte) hash.Hash
	algorithm string
	sz        ObjectCodec
}

type Header struct {
	Algorithm   string `json:"alg,omitempty"`
	ContentType string `json:"cty,omitempty"`
	KeyID       string `json:"kid,omitempty"`
	Type        string `json:"typ,omitempty"`
}

// Payload is a JWT payload according to the RFC 7519.
type Payload struct {
	Issuer         string     `json:"iss,omitempty"`
	Subject        string     `json:"sub,omitempty"`
	Audience       Audience   `json:"aud,omitempty"`
	ExpirationTime *core.Time `json:"exp,omitempty"`
	NotBefore      *core.Time `json:"nbf,omitempty"`
	IssuedAt       *core.Time `json:"iat,omitempty"`
	JWTID          string     `json:"jti,omitempty"`
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
			if val, ok := vv[i].(string); ok {
				aud[i] = val
			}
		}
		*a = aud
	}
	return nil
}

func NewJwtCodec(alg string) CryptoCodec {
	j := &jwtCodec{
		sz: JSONEncoder{},
	}

	hashFunc, err := core.JwtHS(alg)
	j.hashFunc = hashFunc
	j.algorithm = alg
	j.err = err
	return j
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
	hb, err := j.sz.Encode(Header{
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
	err = j.sz.Decode(hp, &header)
	if err != nil {
		return nil, err
	}

	hs, err := core.JwtHS(header.Algorithm)
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
