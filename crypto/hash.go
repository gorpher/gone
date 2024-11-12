package crypto

import (
	"crypto/hmac"
	"crypto/md5" // #nosec
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

/*

hash摘要算法常用作：文件校验、数字签名和鉴权协议
常用的hash摘要算法有：
1. md4: 适用于32位处理器。
2. md5: md4升级版本，性能一般，安全性可靠。
3. sha1 sha224 sha256 sha384: 安全性更高，可以理解为md5的升级版。
4. hmac :哈希运算消息认证码（Hash-based Message Authentication Code），需要一个密钥进行签名，也可以理解为对原有hash进行加盐处理，只是这个盐就是要签名的密钥。
hmac可以与基本hash算法进行组合：hmac_sha1 hmac_sha224 hmac_sha256 hmac_sha384 hmac_sha512 hmac_md5 等等...
*/

// MD5 md5摘要算法,可以加盐,参数salt是多参数,实现多态的行为.
func MD5(str []byte, salt ...[]byte) string {
	m5 := md5.New() // nolint
	m5.Write(str)
	if len(salt) > 0 {
		m5.Write(salt[0])
	}
	return hex.EncodeToString(m5.Sum(nil))
}

// HMacMD5 md5摘要算法,使用hmac算法生成hash,参数salt是多参数,实现多态的行为.
// 参考github.com/dxvgef/gommon/encrypt/sha.go中的hash算法.
func HMacMD5(data []byte, salt ...[]byte) string {
	var s []byte
	if len(salt) > 0 {
		s = salt[0]
	}
	var h hash.Hash
	if len(s) > 0 {
		h = hmac.New(md5.New, s)
	} else {
		h = md5.New() // nolint
	}
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// HMacSha256 生成hmac256摘要
func HMacSha256(key []byte, value []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(value)
	return hex.EncodeToString(h.Sum(nil))
}
