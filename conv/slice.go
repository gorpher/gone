package conv

import (
	"strconv"
	"strings"
)

// RemoveSliceValue 删除切片中指定值的元素
func RemoveSliceValue(s []int64, value int64) []int64 {
	index := 0
	for _, v := range s {
		if v != value {
			s[index] = v
			index++
		}
	}
	return s[:index]
}

func Strings2Int64(arrays []string) []int64 {
	result := make([]int64, len(arrays))
	for i := range arrays {
		v, _ := strconv.ParseInt(arrays[i], 10, 64)
		result[i] = v
	}
	return result
}

func Int64sToString(arrays []int64) []string {
	result := make([]string, len(arrays))
	for i := range arrays {
		v := strconv.FormatInt(arrays[i], 10)
		result[i] = v
	}
	return result
}

// InSlice reports whether v is present in elems.
func InSlice[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func InSliceIgnoreCase(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
		if strings.ToLower(v) == strings.ToLower(s) {
			return true
		}
	}
	return false
}

// Contains reports whether v is present in elems.
func Contains[T comparable](elems []T, v T) bool {
	return InSlice(elems, v)
}
