package gone

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// AesEncryptCBC 加密 AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
func AesEncryptCBC(origData, key, iv []byte) (encrypted []byte, err error) {
	var block cipher.Block
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()                 // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)   // 补全码
	blockMode := cipher.NewCBCEncrypter(block, iv) // 加密模式
	encrypted = make([]byte, len(origData))        // 创建数组
	blockMode.CryptBlocks(encrypted, origData)     // 加密
	return encrypted, err
}

// AesDecryptCBC cbc模式解密
func AesDecryptCBC(encrypted []byte, key, iv []byte) (decrypted []byte, err error) {
	var block cipher.Block
	block, err = aes.NewCipher(key) // 分组秘钥
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv) // 加密模式
	decrypted = make([]byte, len(encrypted))       // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)    // 解密
	decrypted = pkcs5UnPadding(decrypted)          // 去除补全码
	return decrypted, err
}

// 填充明文
func pkcs5Padding(plainText []byte, blockSize int) []byte {
	padding := blockSize - len(plainText)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plainText, padtext...)
}

// 去除填充
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
