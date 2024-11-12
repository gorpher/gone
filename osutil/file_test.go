package osutil

import (
	"testing"
)

func TestFileSize(t *testing.T) {
	filename := "./file_test.go"
	size := FileSize(filename)
	if size <= 0 {
		t.Error("获取文件大小错误")
	}
}

func TestFileExist(t *testing.T) {
	if !FileExist(".") {
		t.Error(". 文件不可能存在")
		return
	}
}

func TestFileIsBinary(t *testing.T) {
	if FileIsBinary("not binary content") {
		t.Error("输入内容不是二进制")
		return
	}
}

func TestFileIsDir(t *testing.T) {
	if !FileIsDir(".") {
		t.Error(". 不可能是文件夹")
		return
	}
}

func TestFileFindPath(t *testing.T) {
	s, err := FileFindPath("file_test.go")
	if err != nil {
		t.Error(err)
		return
	}
	if s == "" {
		t.Error("file_test.go文件存在")
	}
}
