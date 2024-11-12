package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"github.com/gorpher/gone/osutil"
	"testing"

	"github.com/tjfoc/gmsm/pkcs12"
	"github.com/tjfoc/gmsm/sm4"
	"golang.org/x/crypto/blowfish" // nolint
	"golang.org/x/crypto/cast5"    // nolint
	"golang.org/x/crypto/tea"      // nolint
	"golang.org/x/crypto/twofish"  // nolint
	"golang.org/x/crypto/xtea"     // nolint
)

func TestAesEncryptCBC(t *testing.T) {
	key := []byte(osutil.RandString(16))
	iv := []byte(osutil.RandString(16))
	txt := []byte(osutil.RandString(1024))
	encrypt, err := EncryptByAesCBC(txt, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	taed, err := DecryptByAesCBC(encrypt, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	if len(txt) != len(taed) {
		t.Fatalf("原文和解密后的原文不匹配 %s!=%s", txt, taed)
	}
	if !bytes.Equal(txt, taed) {
		t.Fatalf("原文和解密后的原文不匹配 %s!=%s", txt, taed)
	}
}

func TestEncryptByAesCTR(t *testing.T) {
	key := []byte(osutil.RandString(16))
	txt := []byte(osutil.RandString(1024))
	encrypt, err := EncryptByAesCTR(txt, key)
	if err != nil {
		t.Fatal(err)
	}
	taed, err := DecryptByAesCTR(encrypt, key)
	if err != nil {
		t.Fatal(err)
	}
	if len(txt) != len(taed) {
		t.Fatalf("原文和解密后的原文不匹配 %s!=%s", txt, taed)
	}
	if !bytes.Equal(txt, taed) {
		t.Fatalf("原文和解密后的原文不匹配 %s!=%s", txt, taed)
	}
}

func aesBlock(key []byte) cipher.Block {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return block
}

func desBlock(key []byte) cipher.Block {
	// des加密只能设置长度为8的字节数组
	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return block
}

func tripleDESBlock(key []byte) cipher.Block {
	// des加密只能设置长度为8的字节数组
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		panic(err)
	}
	return block
}

// 4. BlowFish算法
// BlowFish算法是一个64位分组及可变密钥长度的分组密码算法，该算法是非专利的。
// BlowFish算法基于Feistel网络（替换/置换网络的典型代表），加密函数迭代执行16轮。分组长度为64位（bit），密钥长度可以从32位到448位。
// 算法由两部分组成：密钥扩展部分和数据加密部分。密钥扩展部分将最长为448位的密钥转化为共4168字节长度的子密钥数组。
// 其中，数据加密由一个16轮的Feistel网络完成。每一轮由一个密钥相关置换和一个密钥与数据相关的替换组成。
func blowfishBlock(key []byte) cipher.Block {
	blowfishBlock, err := blowfish.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return blowfishBlock
}

func cast5Block(key []byte) cipher.Block {
	block, err := cast5.NewCipher(key)
	if err != nil {
		panic(err)
	}

	return block
}

func twofishBlock(key []byte) cipher.Block {
	block, err := twofish.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return block
}

// TEA（Tiny Encryption Algorithm）算法。
// 分组长度为64位，密钥长度为128位。采用Feistel网络。
// 其作者推荐使用32次循环加密，即64轮。
// TEA算法简单易懂，容易实现。但存在很大的缺陷，如相关密钥攻击。
// 由此提出一些改进算法，如XTEA。
func teaBlock(key []byte) cipher.Block {
	block, err := tea.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return block
}

func xteaBlock(key []byte) cipher.Block {
	block, err := xtea.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return block
}

func sm4Block(key []byte) cipher.Block {
	block, err := sm4.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return block
}

func pkcs12Block(key []byte) cipher.Block {
	block, err := pkcs12.New(key, len(key))
	if err != nil {
		panic(err)
	}
	return block
}

func TestBlockEncrypt(t *testing.T) {
	txt := []byte(osutil.RandString(1024))
	var tests = []struct {
		newBlock func(key []byte) cipher.Block
		stream   BlockStreamMode
		keyLen   int
	}{
		{desBlock, CTR, 8},
		{desBlock, OFB, 8},
		{desBlock, CFB, 8},
		{desBlock, RC4, 8},

		{tripleDESBlock, CTR, 24},
		{tripleDESBlock, OFB, 24},
		{tripleDESBlock, CFB, 24},
		{tripleDESBlock, RC4, 24},

		{aesBlock, CTR, 16},
		{aesBlock, OFB, 16},
		{aesBlock, CFB, 16},
		{aesBlock, RC4, 16},

		{blowfishBlock, CTR, 32},
		{blowfishBlock, OFB, 32},
		{blowfishBlock, CFB, 32},
		{blowfishBlock, RC4, 32},

		{cast5Block, CTR, 16},
		{cast5Block, OFB, 16},
		{cast5Block, CFB, 16},
		{cast5Block, RC4, 16},

		{twofishBlock, CTR, 16},
		{twofishBlock, OFB, 16},
		{twofishBlock, CFB, 16},
		{twofishBlock, RC4, 16},

		{teaBlock, CTR, 16},
		{teaBlock, OFB, 16},
		{teaBlock, CFB, 16},
		{teaBlock, RC4, 16},

		{xteaBlock, CTR, 16},
		{xteaBlock, OFB, 16},
		{xteaBlock, CFB, 16},
		{xteaBlock, RC4, 16},

		{sm4Block, CTR, 16},
		{sm4Block, OFB, 16},
		{sm4Block, CFB, 16},
		{sm4Block, RC4, 16},

		{pkcs12Block, CTR, 16},
		{pkcs12Block, OFB, 16},
		{pkcs12Block, CFB, 16},
		{pkcs12Block, RC4, 16},
	}
	for i := range tests {
		key := []byte(osutil.RandString(tests[i].keyLen))
		encrypt, err := BlockEncrypt(tests[i].newBlock(key), tests[i].stream, txt)
		if err != nil {
			t.Fatal(err)
		}
		taed, err := BlockDecrypt(tests[i].newBlock(key), tests[i].stream, encrypt)
		if err != nil {
			t.Fatal(err)
		}
		if len(txt) != len(taed) {
			t.Fatalf("原文和解密后的原文不匹配 %s!=%s", txt, taed)
		}
		if !bytes.Equal(txt, taed) {
			t.Fatalf("原文和解密后的原文不匹配 %s!=%s", txt, taed)
		}
	}
}
