package gone

import (
	"bytes"
	"strings"
	"unsafe"
)

// BytesToStr []byte转string
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StrToBytes string转[]byte
func StrToBytes(str string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&str))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// Contains 判断字符串是否存在切片中
func Contains(str string, s []string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// StringReplaceIgnoreCase 忽略大小写替换字符串
func StringReplaceIgnoreCase(text, source, target string) string {
	buf := &bytes.Buffer{}
	textLower := strings.ToLower(text)
	searchStrLower := strings.ToLower(source)
	searchStrLen := len(source)
	var end int
	for {
		idx := strings.Index(textLower, searchStrLower)
		if 0 > idx {
			break
		}
		buf.WriteString(text[:idx])
		buf.WriteString(target)
		end = idx + searchStrLen
		textLower = textLower[end:]
	}
	buf.WriteString(text[end:])
	return buf.String()
}
