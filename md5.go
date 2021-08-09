package gone

import (
	"crypto/hmac"
	"crypto/md5" // #nosec
	"encoding/hex"
	"hash"
)

// MD5 md5摘要算法,可以加盐,参数salt是多参数,实现多态的行为.
func MD5(str []byte, salt ...[]byte) (string, error) {
	var err error
	m5 := md5.New()
	_, err = m5.Write(str)
	if err != nil {
		return "", err
	}
	if len(salt) > 0 {
		_, err = m5.Write(salt[0])
		if err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(m5.Sum(nil)), nil
}

// MD5Encrypt md5摘要算法,使用hmac算法生成hash,参数salt是多参数,实现多态的行为.
// 参考github.com/dxvgef/gommon/encrypt/sha.go中的hash算法.
func MD5Encrypt(data []byte, salt ...[]byte) (string, error) {
	var s []byte
	var err error
	if len(salt) > 0 {
		s = salt[0]
	}
	var h hash.Hash
	if len(s) > 0 {
		h = hmac.New(md5.New, s)
	} else {
		h = md5.New() // #gosec
	}
	_, err = h.Write(data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
