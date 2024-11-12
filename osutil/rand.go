package osutil

import (
	"bytes"
	rnd "crypto/rand"
	"log"
	"math/rand"
	"time"
)

// RandInt64 指定范围内的随机数字，max必须大于min.
func RandInt64(min, max int64) int64 {
	rand.Seed(rand.Int63n(time.Now().UnixNano())) //nolint
	return min + rand.Int63n(max-min)
}

// RandInt32 指定范围内的随机数字，max必须大于min.
func RandInt32(min, max int32) int32 {
	rand.Seed(rand.Int63n(time.Now().UnixNano())) //nolint
	return min + rand.Int31n(max-min)
}

// RandInt 指定范围内的随机数字.
func RandInt(min, max int) int {
	rand.Seed(rand.Int63n(time.Now().UnixNano())) //nolint
	return min + rand.Intn(max-min)
}

// RandInts 生成指定范围int类型数组.
func RandInts(from, to, size int) []int {
	if to-from < size {
		size = to - from
	}

	var slice []int
	for i := from; i < to; i++ {
		slice = append(slice, i)
	}

	var ret []int
	for i := 0; i < size; i++ {
		idx := rand.Intn(len(slice))
		ret = append(ret, slice[idx])
		slice = append(slice[:idx], slice[idx+1:]...)
	}
	return ret
}

// =======================================================================================================================

// RandLower 指定长度的随机小写字母.
func RandLower(l int) string {
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		if string(RandInt32(97, 122)) != temp {
			temp = string(RandInt32(97, 122))
			if _, err := result.WriteString(temp); err != nil {
				return ""
			}
			i++
		}
	}
	return result.String()
}

// RandUpper 指定长度的随机大写字母.
func RandUpper(l int) string {
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		if string(RandInt32(65, 90)) != temp {
			temp = string(RandInt32(65, 90))
			if _, err := result.WriteString(temp); err != nil {
				return ""
			}
			i++
		}
	}
	return result.String()
}

// ======================================================================================================================

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

// RandBytes 生成随机长度字节.
func RandBytes(length int) []byte {
	randomBytes := make([]byte, length)
	var err error
	_, err = rnd.Read(randomBytes)
	if err != nil {
		log.Fatal("Unable to generate random bytes")
	}
	return randomBytes
}

// RandAlphaString 生成随机长度字母.
// Deprecated
func RandAlphaString(length int) string {
	result := make([]byte, length)
	bufferSize := int(float64(length) * 1.3)
	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			randomBytes = RandBytes(bufferSize)
		}
		if idx := int(randomBytes[j%length] & letterIdxMask); idx < len(letterBytes) {
			result[i] = letterBytes[idx]
			i++
		}
	}
	return string(result)
}

// RandString 生成随机长度字符串,推荐使用.
func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
