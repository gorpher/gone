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

// Base64RawStdEncode baseRaw标准编码
func Base64RawStdEncode(value []byte) []byte {
	encoded := make([]byte, base64.RawStdEncoding.EncodedLen(len(value)))
	base64.RawStdEncoding.Encode(encoded, value)
	return encoded
}

// Base64RawStdDecode base标准解码
func Base64RawStdDecode(value []byte) ([]byte, error) {
	decoded := make([]byte, base64.RawStdEncoding.DecodedLen(len(value)))
	b, err := base64.RawStdEncoding.Decode(decoded, value)
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

// Base64RawURLEncode baseRawURL编码
func Base64RawURLEncode(value []byte) []byte {
	encoded := make([]byte, base64.RawURLEncoding.EncodedLen(len(value)))
	base64.RawURLEncoding.Encode(encoded, value)
	return encoded
}

// Base64RawURLDecode baseRawURL解码
func Base64RawURLDecode(value []byte) ([]byte, error) {
	decoded := make([]byte, base64.RawURLEncoding.DecodedLen(len(value)))
	b, err := base64.RawURLEncoding.Decode(decoded, value)
	if err != nil {
		return nil, err
	}
	return decoded[:b], nil
}
