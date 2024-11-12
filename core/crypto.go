package core

type (
	SecretKeyLengthType int
	SecretKeyFormatType int
)

const (
	PKCS8 SecretKeyFormatType = 1
	PKCS1 SecretKeyFormatType = 2
)

const (
	RSA SecretKeyLengthType = 1
	M2  SecretKeyLengthType = 2
)

func SecretKeyLengthTypeParse(v int) SecretKeyLengthType {
	switch v {
	case 1:
		return RSA
	case 2:
		return M2
	}
	return RSA
}

func SecretKeyFormatTypeParse(v int) SecretKeyFormatType {
	switch v {
	case 1:
		return PKCS8
	case 2:
		return PKCS1
	}
	return PKCS8
}
