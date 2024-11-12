package osutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileSize 获取文件大小.
func FileSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return fi.Size()
}

// FileExist 判断文件是否存在.
func FileExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

// FileIsBinary 判断文件是否是二进制文件.
func FileIsBinary(content string) bool {
	for _, b := range content {
		if b == 0 {
			return true
		}
	}

	return false
}

// FileIsDir 判断文件是否是目录.
func FileIsDir(path string) bool {
	fio, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}

	return fio.IsDir()
}

// FileCopy 复制文件，代码仅供参考.
// Deprecated
func FileCopy(source, dest string) (err error) {
	var sourcefile *os.File
	sourcefile, err = os.Open(filepath.Clean(source))
	if err != nil {
		return err
	}

	defer sourcefile.Close() //nolint

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close() //nolint

	_, err = io.Copy(destfile, sourcefile)
	if err != nil {
		return err
	}
	if sourceinfo, errv := os.Stat(source); errv != nil {
		err = os.Chmod(dest, sourceinfo.Mode())
	}
	return err
}

// FileMove 移动文件，代码仅供参考.
// Deprecated
func FileMove(src, dest string) error {
	dir := filepath.Dir(dest)
	_, err := os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, os.FileMode(0755))
		if err != nil {
			return err
		}
	}
	return os.Rename(src, dest)
}

// FileFindPath 获取文件名路径，首先判断文件是否可以直接访问，优先获取当前可执行文件夹下，再去找工作路径下.
func FileFindPath(fname string) (string, error) {
	if filepath.IsAbs(fname) && !FileExist(fname) {
		return "", fmt.Errorf("config file %s not exist", fname)
	}
	location, err := os.Executable()
	if err != nil {
		return "", err
	}
	fname = filepath.Join(filepath.Dir(location), fname)
	if !FileExist(fname) {
		location, err := os.Getwd()
		if err != nil {
			return "", err
		}
		_, fname = filepath.Split(fname)
		fname = filepath.Join(location, fname)
		if !FileExist(fname) {
			return "", fmt.Errorf("config file %s not exist", fname)
		}
	}
	return filepath.Clean(fname), nil
}
