package gone

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 md5摘要算法.
func MD5(str []byte) string {
	m := md5.New()
	m.Write(str)
	return hex.EncodeToString(m.Sum(nil))
}

// MD5Encrypt md5摘要算法加盐.
func MD5Encrypt(str, salt []byte) string {
	m5 := md5.New()
	m5.Write(str)
	m5.Write(salt)
	return hex.EncodeToString(m5.Sum(nil))
}
