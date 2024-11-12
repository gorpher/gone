package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gorpher/gone/core"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

var licenseStr = "需要签名的字符串"

func TestGenerateBase64Key(t *testing.T) {
	pkStr, pbkStr, err := GenerateBase64Key(core.M2, core.PKCS1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("私钥：", pkStr)
	t.Log("公钥：", pbkStr)

	pkStr, pbkStr, err = GenerateBase64Key(core.M2, core.PKCS8)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("私钥：", pkStr)
	t.Log("公钥：", pbkStr)

	pkStr, pbkStr, err = GenerateBase64Key(core.RSA, core.PKCS1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("私钥：", pkStr)
	t.Log("公钥：", pbkStr)

	pkStr, pbkStr, err = GenerateBase64Key(core.RSA, core.PKCS8)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("私钥：", pkStr)
	t.Log("公钥：", pbkStr)
}

const (
	PrivateKeyPemStr = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgH1D0jXGNiOvh0M8tBC5x+0ZSQLnG2yfii+xkHm4PrAWpQ2rS3bT
HbU2fPmDw0Aaxs4NWGCoWWrMcht71myCNqFAwDM/CtVZdhKEFclEUdJCKuL4qwJ4
e60JJLCS2QUEy/n6bA1WovXBTeR8tSI+OKUm9EvFmNGQxhmlSfLNpw/1AgMBAAEC
gYAOLi5Ozhh047sBPo730cAzNAiS3oy5ODpRed1sGhJmprmamYiac/3J9Ngi+uqQ
iDd3PgWCM6yjrW9Bczxr3jXG07HqB2XdcJvXlpFpZNmtujHrbKEkBDMbWT3DH5Hq
P/7MOWmcmoMnUqGlL7dXZ6EL/zn3Asy3hGizwYqXyCzXIQJBAOOQLk9zyVQlj0Yn
KCcUWJUytYK//SrKDbGnjjnCROX+TwgjcWZ+MRVx0wTBspEXh4GH0OhR6e2jVI50
+5k7Ks0CQQCM6xX4mYl33Vv+lEJrtEtE+wUsVO6NYz2KLJ59Nx+RJ5UzTHBRdXvE
f36zu0pyAUeYlr3ZP8khLbyyIaB29knJAkEAoeffSyAySfA/M8aARu2u6NgfVFuM
oHkJrTBtfKK/qnN5f2zYLffyrDND08qMZba77mjXNbOyICVo78JDkA4MsQJAH65i
nCd4nngnzI5seGZqXbHJsfPORf8/wKbTYvdXo3ywsH3I6qdtEfpP8/xxejwLaqTJ
PeR3RXxQ5gNlXhl08QJAD9x2awK64KR73dG/PjUtph9dp88FoeVgsrjdFSyvEqaA
VIHgpDIFFTxcDGQ3U5qdln/ekvRIF5WEXD6J2GITMg==
-----END RSA PRIVATE KEY-----`

	PublicKeyPemStr = `-----BEGIN PUBLIC KEY-----
MIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgH1D0jXGNiOvh0M8tBC5x+0ZSQLn
G2yfii+xkHm4PrAWpQ2rS3bTHbU2fPmDw0Aaxs4NWGCoWWrMcht71myCNqFAwDM/
CtVZdhKEFclEUdJCKuL4qwJ4e60JJLCS2QUEy/n6bA1WovXBTeR8tSI+OKUm9EvF
mNGQxhmlSfLNpw/1AgMBAAE=
-----END PUBLIC KEY-----`

	PrivateKeyBase64Str = `MIICXgIBAAKBgQC+6sp7l7dmxZ2vExLZFIkekhs9gV1E0d+9xXhNcReXhcOMFCYBdVLgDdJtHddNsu5CqBBZbBOkDl4lPV7u1AqbS1AnENuXJaAaaCDutpfxZZUoUyi4ZcbtcbtU/3g0eblosE5xxUUpq30JIrhz83biq8tVE1Wbx3cZMvgQy5QHbQIDAQABAoGBAI35TFpcmKZ0jq6DIKEOBGoXfOpgKVvkNt6I2s28LC8h6ilhUmIDPX4gyTsb1eCSD1zCXmYhWPnHNXu8B7zTMo6/C3OxNFYmLnWzm4EojkL4K50jDusJp/eWrSDp4Gg/lTxnZPYE/Q5138fF98BVv7G0hCK8dowKHJdfHdrKE1wBAkEA9OKfL62Te+fohALoUsQVzk1idVwWMMzlYhlXW4+nAvI/s/DFvS50Lw2l/WJmA+MEUrfOeCj8hEjOkttk1PS/pQJBAMeVGfd+UoaQjCTozLbIXf6UJEFpuRsowPdnNCYAqYjx4DiI4nC1SHb7P/woIgkghSeLFU+RxGgdeYt+XRZ8HikCQQDi2jkLKunQQ8JS4HqliX6F0YwfGgJ4jKcGHGGfsVDO2ukGYUpc+XapzCPzub61ZQ0xL5L2H0nlpaivxMwAtwX9AkAug7d7kPtW1VV0PLWJXAVcEdapUCSOCd9/SZRDzx+0BPtG8dAkiHuND12IPSpBiky+PJII62YlBcmQEzFKzj6RAkEApvOBh7AeRSGPRxo8SiSO93mJnYx6woC8869BH5xIhv5ZdDmc4KzSZjijkfZ216qQelHbTwewU9yykbP9++ekhA==`
	PublicKeyBase64Str  = `MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC+6sp7l7dmxZ2vExLZFIkekhs9gV1E0d+9xXhNcReXhcOMFCYBdVLgDdJtHddNsu5CqBBZbBOkDl4lPV7u1AqbS1AnENuXJaAaaCDutpfxZZUoUyi4ZcbtcbtU/3g0eblosE5xxUUpq30JIrhz83biq8tVE1Wbx3cZMvgQy5QHbQIDAQAB`

	PrivateKeyHexStr = `308204a50201000282010100dad2af9734a1b1ec6683f1352e607b3114da7211309a9802eb20ac4bab6fcf404c96042766e8d6efd36e18b79e004f77e26c36cbc5b7e27922ec81ceedc8b714c565b748043689fd362038643e4167bac83feeb603de18e24dc3f66095d961dcc151c9ecf119bf80523d347ffe788f0b0cac4a021c19e8ccdbd1213680c7c93a489fdbbc176ee39c0ca76d714add28e64da01536372b03bcb054575589283c2e33c164e8f73610572568606b6f2134d868756311d8641bda47ab96d2ab9dcd25e1a52890c15bb820356a15ce5e4669b134ae2473c8d038de9515c114f4ef1613b959f6945be07664caa9c99eac98c9310d6d30a600f4108237ebbd456e1f112d02030100010282010100bad0ea65ede229f3886616b4ef7e214e7ade304a2ab3a119c4c16537490ab0e6d53bb21a2d8a958db4752717040681d1f8f5a8267a0b8e871ae0cdc5eb4dd3b820fac0e9f3e6d811dc76cf8e0d746b699472b88a9e6cabd3f0ce5f76801851ab55444d5f434b5729e78c27592ad8a44eaba81f9b0380bf36be6821b1d56a3b8939a4a18ea9eb165ff41e6f5be1da4bcb2f83fa8df4ffd52dea86b0cfb214fb6102603bf79b4adf88abbf1b4323a2b21337d314b8e9c1170747f914cbf87f6cc4f12e89447f28f90bfa536bc3395b258e95eabb841d64934cee735d5160ca43670f31102de03f9b4b250c1cc1607d8e51433f2b8b5e2cf57f01e25a822c5aede502818100f1b39dca82c47fabf4dd8e54f8ba567c59d3092c24545c730a5b0470a577e186109a382d8bb8829aa44d8062cf24a59dbc69b72db6223f946c8f4b4ff770ca467eb3eabc706139a7e1d627febfd2f4e9aa9bfda5f1cb408fe376a132e68a81d1e2461d60fe1b2b89ae0073138ed841610b48d52845f1b3706ec0295dd712b7b302818100e7c4972197d792496be0ef20b52ccd74dbd7df473795b1fb4b9fe9eb0ebd97e3cc28c2c6652911b0770ef012f3c11b8ede9efd1c13505b3d4a6a57ff264ee41d681d6efe02a197a99ed7857aa38f6037fba538e7b93aedeebd2fd9e14e2be5d5ee312a864ff5b633857f46ea847630952f1ca9975f5ca111d7f3a8db9358a39f0281805a70f68b4f914da0bf98a3e8c1c5a01519db70e43697e69c1974e35d6f5d43635215130e5fe8e3de0fbafc5e7cda5eaa7e552479135d0f636f97d2fb92407f400fab2d1be4054d78b775d63369fdfb2cf06d3c657aebae35e94c7b973b52faaed9b798c8b16ce346ba786a9717ed6dd16d528c886c5bbbe4475cda5dc5dbb82702818100d4ad9eedd1b39ce6c91ad8f48facb440b6f86a48a4e63633de9ab901dd3df7a2af16fc5d28393ea54b2ba6fc0d38383cab6703e6fe862fa397a4ec6913d3331b150e656aac2972cdd117fec1a253903cef2c1782f483f210b104b7103c36a62ae0efb7111750e7c87189711f053c9baa5a5817fbf323421ee8a70c5da9e19e0b02818100d1c9e36e4c56ad882a9ff59aa796983fa88921dfcd285a19c99fe9040bae7e1ab78b8197fb037db97c8f338526824ea663f00e041926c0a24b8b8091820f254701f32b8701c99d6f33621ad0551d191ceaf31ffe7549e1feb7ec47a487eced56d3de15bfb603cca6efcd95f810dc166688f971bf3ed9f4ada1adb148d9f4f4c4`
	PublicKeyHexStr  = `3082010a0282010100dad2af9734a1b1ec6683f1352e607b3114da7211309a9802eb20ac4bab6fcf404c96042766e8d6efd36e18b79e004f77e26c36cbc5b7e27922ec81ceedc8b714c565b748043689fd362038643e4167bac83feeb603de18e24dc3f66095d961dcc151c9ecf119bf80523d347ffe788f0b0cac4a021c19e8ccdbd1213680c7c93a489fdbbc176ee39c0ca76d714add28e64da01536372b03bcb054575589283c2e33c164e8f73610572568606b6f2134d868756311d8641bda47ab96d2ab9dcd25e1a52890c15bb820356a15ce5e4669b134ae2473c8d038de9515c114f4ef1613b959f6945be07664caa9c99eac98c9310d6d30a600f4108237ebbd456e1f112d0203010001`
)

func TestDecodePemHexBase64(t *testing.T) {
	_, err := DecodePemHexBase64([]byte(PrivateKeyPemStr))
	if err != nil {
		t.Fatal(err)
	}
	_, err = DecodePemHexBase64([]byte(PrivateKeyBase64Str))
	if err != nil {
		t.Fatal(err)
	}
	_, err = DecodePemHexBase64([]byte(PrivateKeyHexStr))
	if err != nil {
		t.Fatal(err)
	}
	_, err = DecodePemHexBase64([]byte(PublicKeyPemStr))
	if err != nil {
		t.Fatal(err)
	}
	_, err = DecodePemHexBase64([]byte(PublicKeyBase64Str))
	if err != nil {
		t.Fatal(err)
	}
	_, err = DecodePemHexBase64([]byte(PublicKeyHexStr))
	if err != nil {
		t.Fatal(err)
	}
}

func testRandStr(l int) string {
	buff := make([]byte, l)
	rand.Read(buff) //nolint
	str := base64.StdEncoding.EncodeToString(buff)
	// Base 64 can be longer than len
	return str[:l]
}

func TestGenerateSSHKey(t *testing.T) {
	var publicKey ssh.PublicKey
	var pkBytes []byte
	var pbkBytes []byte
	var err error
	pkBytes, pbkBytes, err = GenerateSSHKey(RSA2048)
	if err != nil {
		t.Fatal(err)
	}
	key, _, _, _, err := ssh.ParseAuthorizedKey(pbkBytes)
	if err != nil {
		t.Fatal(err)
	}
	if key == nil {
		t.Fatal("key not found")
	}
	// 更多的例子，请查看。
	// golang.org/x/crypto/ssh/testdata_test.go:26
	// golang.org/x/crypto/ssh//keys_test.go:577
	publicKey, err = ssh.ParsePublicKey(key.Marshal())
	if err != nil {
		t.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(pkBytes)
	if err != nil {
		t.Fatal(err)
	}
	data := []byte("some message")
	signature, err := signer.Sign(rand.Reader, data)
	if err != nil {
		t.Fatal(err)
	}
	err = publicKey.Verify(data, signature)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateECDSAKeyToMemory(t *testing.T) {
	privateBytes, publicBytes, err := GenerateECDSAKeyToMemory(elliptic.P256())
	if err != nil {
		t.Fatal(err)
	}
	block, _ := pem.Decode(privateBytes)
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	pbkBlock, _ := pem.Decode(publicBytes)
	pbk, err := x509.ParsePKIXPublicKey(pbkBlock.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, ok := pbk.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal(errors.New("ecdsa publicKey error"))
	}

	text := "hello world"
	randSign := "12345678901234567"
	hashText := sha256.Sum256([]byte(text))
	r, s, err := ecdsa.Sign(strings.NewReader(randSign), privateKey, hashText[:])
	if err != nil {
		t.Fatal(err)
	}
	b := ecdsa.Verify(publicKey, hashText[:], r, s)
	fmt.Println(b)
}
