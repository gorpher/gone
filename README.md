### gone 
gone是一个简单、方便、小巧好用的工具包.

[![GoDoc](https://pkg.go.dev/badge/github.com/gorpher/gone)](https://pkg.go.dev/github.com/gorpher/gone)
[![Build Status](https://api.travis-ci.com/gorpher/gone.svg?branch=main&status=passed)](https://travis-ci.org/gorpher/gorpher)

#### 非对称加密解密与签名
```go
// GenerateBase64Key 生成base64编码的公私钥
func GenerateBase64Key(secretLength SecretKeyLengthType, secretFormat SecretKeyFormatType) (pkStr string, pbkStr string, err error)

// SignBySM2Bytes 使用sm2私钥签名字符串，返回base64编码的license
func SignBySM2Bytes(privateKey, licenseBytes []byte) (license string, err error) 

// SignBySM2  使用sm2私钥对象指针签名字符串，返回base64编码的license
func SignBySM2(privateKey *sm2.PrivateKey, licenseBytes []byte) (license string, err error)

// SignByRSABytes 使用rsa私钥签名字符串，返回base64编码的license
func SignByRSABytes(key, licenseBytes []byte) (string, error) 

// SignByRSA 使用rsa私钥对象指针签名字符串，返回base64编码的license
func SignByRSA(key *rsa.PrivateKey, licenseBytes []byte) (license string, err error)

// VerifyBySM2 使用sm2公钥验证签名的license
func VerifyBySM2(publicKeyBase64, licenseCode string) (license string, valid bool, err error)

// VerifyByRSA 使用rsa公钥验证签名的license
func VerifyByRSA(publicKeyBase64, licenseCode string) (license string, valid bool, err error)

// RsaPublicEncrypt Rsa公钥加密，参数publicKeyStr必须是hex、base64或者是pem编码
func RsaPublicEncrypt(publicKeyStr string, textBytes []byte) ([]byte, error) 

// ParsePublicKey 解析公钥，derBytes可以使用DecodePemHexBase64函数获取
func ParsePublicKey(derBytes []byte) (publicKey *rsa.PublicKey, err error)

// RsaPrivateDecrypt 解析rsa私钥，参数privateKeyStr必须是hex、base64或者是pem编码
func RsaPrivateDecrypt(privateKeyStr string, cipherBytes []byte) (textBytes []byte, err error)

// ParsePrivateKey 解析私钥，derBytes可以使用DecodePemHexBase64函数获取
func ParsePrivateKey(derBytes []byte) (privateKey *rsa.PrivateKey, err error) 

// DecodePemHexBase64 解析pem或者hex或者base64编码成der编码
func DecodePemHexBase64(keyStr string) ([]byte, error)

```


### 随机函数
```go
// RandInt64 指定范围内的随机数字，max必须大于min。
func RandInt64(min int64, max int64) int64

// RandInt32 指定范围内的随机数字，max必须大于min。
func RandInt32(min int32, max int32) int32 

// RandInt 指定范围内的随机数字
func RandInt(min int, max int) int 

// RandInts 生成指定范围int类型数组
func RandInts(from, to, size int) []int

// RandLower 指定长度的随机小写字母
func RandLower(l int) string

// RandUpper 指定长度的随机大写字母
func RandUpper(l int) string


// RandBytes 生成随机长度字节
func RandBytes(length int) []byte 

// RandString 生成随机长度字符串,推荐使用。
func RandString(n int) string 
```


### 文件处理函数
```go
// FileSize 获取文件大小
func FileSize(path string) 

// FileSize 判断文件是否存在
func FileExist(file string) bool

// FileIsBinary 判断文件是否是二进制文件
func FileIsBinary(content string) bool

// FileIsDir 判断文件是否是目录
func FileIsDir(path string) bool

// FileFindAPath 获取文件名路径，首先判断文件是否可以直接访问，优先获取当前可执行文件夹下，再去找工作路径下。
func FileFindPath(fname string)

```

### aes加解密函数
```go
// AesEncryptCBC 加密 AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
func AesEncryptCBC(origData, key, iv []byte) (encrypted []byte, err error)

// AesDecryptCBC cbc模式解密
func AesDecryptCBC(encrypted []byte, key, iv []byte) (decrypted []byte, err error)
```

### 字符串转换函数
```go
// BytesToStr []byte转string
func BytesToStr(b []byte) string

// StrToBytes string转[]byte
func StrToBytes(str string) []byte

// Contains 判断字符串是否存在切片中
func Contains(str string, s []string) bool

// StringReplaceIgnoreCase 忽略大小写替换字符串
func StringReplaceIgnoreCase(text, source, target string) string

```

### 时间格式化函数
```go
// TimeToStr 返回时间的字符串格式
func TimeToStr(t time.Time, format ...string) string 

// Timestamp将unix时间转为时间字符串
func TimestampToStr(t int64, format ...string) string

// FormatByStr 将字符串中的时间变量（y年/m月/d日/h时/i分/s秒）转换成时间字符串
func FormatByStr(tpl string, t int64) string 

// GetMonthRange 获得指定年份和月份的起始unix时间和截止unix时间
func GetMonthRange(year int, month int) (beginTime, endTime int64, err error) 

// GetWeek 获得星期的数字
func GetWeek(t time.Time) int


```

### 其他函数
```go
//MacAddr 获取机器mac地址，返回mac字串数组
func MacAddr() (upMac []string, err error)

```
