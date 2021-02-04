package gone

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// FileSize 获取文件大小
func FileSize(path string) int64 {
	fi, err := os.Stat(path)
	if nil != err {
		return 0
	}

	return fi.Size()
}

// FileSize 判断文件是否存在
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

// FileIsBinary 判断文件是否是二进制文件
func FileIsBinary(content string) bool {
	for _, b := range content {
		if 0 == b {
			return true
		}
	}

	return false
}

// FileIsDir 判断文件是否是目录
func FileIsDir(path string) bool {
	fio, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return false
	}
	if nil != err {
		return false
	}

	return fio.IsDir()
}

// FileCopy 复制文件，代码仅供参考
// Deprecated
func FileCopy(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}
	return nil
}

// FileMove 移动文件，代码仅供参考
// Deprecated
func FileMove(src, dest string) error {
	dir := filepath.Dir(dest)
	_, err := os.Stat(dir)
	if err != nil {
		err := os.MkdirAll(dir, os.FileMode(0755))
		if err != nil {
			return err
		}
	}
	return os.Rename(src, dest)
}

// FileFindAPath 获取文件名路径，首先判断文件是否可以直接访问，优先获取当前可执行文件夹下，再去找工作路径下。
func FileFindPath(fname string) (string, error) {
	if filepath.IsAbs(fname) {
		if !FileExist(fname) {
			return "", fmt.Errorf("配置文件%s不存在", fname)
		}
	} else {
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
				return "", fmt.Errorf("配置文件%s不存在", fname)
			}
		}

	}
	return filepath.Clean(fname), nil
}
