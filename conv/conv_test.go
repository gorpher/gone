package conv

import "testing"

func TestContains(t *testing.T) {
	a := "a"
	arr := []string{"a", "b"}
	if !Contains(arr, a) {
		t.Errorf("%s 不在 %#v中", a, arr)
		return
	}
}

func TestStrToBytes(t *testing.T) {
	s := "不是所有人的需求都应该被满足"
	bytes := StrToBytes(s)
	if s2 := BytesToStr(bytes); s != s2 {
		t.Errorf("字符串转换成字节失败： %s!=%s", s, s2)
	}
}
