package gone

import "encoding/base64"

// Base64StdEncode base标准编码
func Base64StdEncode(value []byte) []byte {
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(value)))
	base64.StdEncoding.Encode(encoded, value)
	return encoded
}

// Base64StdDecode base标准解码
func Base64StdDecode(value []byte) ([]byte, error) {
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(value)))
	b, err := base64.StdEncoding.Decode(decoded, value)
	if err != nil {
		return nil, err
	}
	return decoded[:b], nil
}

// Base64URLEncode baseURL编码
func Base64URLEncode(value []byte) []byte {
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(value)))
	base64.URLEncoding.Encode(encoded, value)
	return encoded
}

// Base64URLDecode baseURL解码
func Base64URLDecode(value []byte) ([]byte, error) {
	decoded := make([]byte, base64.URLEncoding.DecodedLen(len(value)))
	b, err := base64.URLEncoding.Decode(decoded, value)
	if err != nil {
		return nil, err
	}
	return decoded[:b], nil
}
