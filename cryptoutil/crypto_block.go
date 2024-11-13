package cryptoutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rc4" // nolint
	"errors"
	"fmt"
	"io"
)

/*
常用五种分组模式：ECB/CBC/CFB/OFB/CTR
1.电码本模式 Electronic Codebook Book (ECB)
2.密码分组链接模式 Cipher Block Chaining (CBC)
3.计算器模式 Counter (CTR)
4.密码反馈模式 Cipher FeedBack (CFB)
5.输出反馈模式 Output FeedBack (OFB)
*/

// EncryptByAesCBC 加密 AES-128 key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256.
// CBC比EBC更安全，但不可并行.
func EncryptByAesCBC(origData, key, iv []byte) (encrypted []byte, err error) {
	var block cipher.Block
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()                 // 获取秘钥块的长度
	origData = pkcs7Padding(origData, blockSize)   // 数据补齐，使用pkcs7填充
	blockMode := cipher.NewCBCEncrypter(block, iv) // 加密模式
	encrypted = make([]byte, len(origData))        // 创建数组
	blockMode.CryptBlocks(encrypted, origData)     // 加密
	return encrypted, err
}

// DecryptByAesCBC cbc模式解密.
func DecryptByAesCBC(encrypted, key, iv []byte) (decrypted []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("DecryptByAesCBC failed")
		}
	}()
	var block cipher.Block
	block, err = aes.NewCipher(key) // 分组秘钥
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv) // 加密模式
	decrypted = make([]byte, len(encrypted))       // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)    // 解密
	decrypted = pkcs7UnPadding(decrypted)          // 去除补全码
	return decrypted, err
}

// pkcs7Padding 填充明文.
func pkcs7Padding(plainText []byte, blockSize int) []byte {
	padding := blockSize - len(plainText)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plainText, padtext...)
}

// pkcs7UnPadding 去除填充.
func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

/*
1. ECB: 这种模式是将整个明文分成若干段相同的小段，然后对每一小段进行加密。
2. CBC: 这种模式是先将明文切分成若干小段，然后每一小段与初始块或者上一段的密文段进行异或运算后，再与密钥进行加密。
3. CTR: 计算器模式不常见，有一个自增的算子，这个算子用密钥加密之后的输出和明文异或的结果得到密文，相当于一次一密。
		  这种加密方式简单快速，安全可靠，而且可以并行加密，但是在计算器不能维持很长的情况下，密钥只能使用一次 。

4. CFB: 密文没有规律, 明文分组是和一个数据流进行的按位异或操作, 最终生成了密文
5. OFB: 密文没有规律, 明文分组是和一个数据流进行的按位异或操作, 最终生成了密文
*/

type BlockStreamMode string

const (
	CTR BlockStreamMode = "CTR"
	CFB BlockStreamMode = "CFB"
	OFB BlockStreamMode = "OFB"
	RC4 BlockStreamMode = "RC4"
)

func BlockEncrypt(block cipher.Block, mode BlockStreamMode, value []byte) ([]byte, error) {
	size := block.BlockSize()
	iv := make([]byte, size, len(value))
	_, err := io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}
	var stream cipher.Stream
	switch mode {
	case CTR:
		stream = cipher.NewCTR(block, iv)
	case CFB:
		stream = cipher.NewCFBEncrypter(block, iv)
	case OFB:
		stream = cipher.NewOFB(block, iv)
	case RC4:
		stream, err = rc4.NewCipher(iv) // nolint
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%s block encrypt unsupported", mode)
	}
	encrypted := make([]byte, len(value))
	stream.XORKeyStream(encrypted, value)
	// 返回偏移量和密文
	iv = append(iv, encrypted...)
	return iv, nil
}

func BlockDecrypt(block cipher.Block, mode BlockStreamMode, value []byte) ([]byte, error) {
	size := block.BlockSize()
	if len(value) > size {
		iv := value[:size]   // 获取偏移量
		value = value[size:] // 提取密文
		var stream cipher.Stream
		var err error
		switch mode {
		case CTR:
			stream = cipher.NewCTR(block, iv)
		case CFB:
			stream = cipher.NewCFBDecrypter(block, iv)
		case OFB:
			stream = cipher.NewOFB(block, iv)
		case RC4:
			stream, err = rc4.NewCipher(iv) // nolint
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("%s block encrypt unsupported", mode)
		}
		stream.XORKeyStream(value, value)
		return value, nil
	}
	return nil, errors.New("decryption failed")
}

func EncryptByAesCTR(origData, key []byte) (encrypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return BlockEncrypt(block, CTR, origData)
}

func DecryptByAesCTR(origData, key []byte) (encrypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return BlockDecrypt(block, CTR, origData)
}
