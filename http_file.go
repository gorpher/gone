package gone

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type HTTPReceiveFile struct {
	fileHeader     *multipart.FileHeader
	MaxReceiveSize int64       //  输入|输出: 最大文件
	DstDirMode     fs.FileMode //  输入|输出: 文件夹名
	DstFileMode    fs.FileMode //  输入|输出: 文件夹名
	DstDir         string      //  输入|输出: 文件夹名
	DstFilename    string      //  输入|输出: 文件名
	FieldName      string      // 输入: 字段名
	FileMIME       string      // 输出: 文件MIME
	TotalSize      int64       // 输出: 总大小
	Hash           hash.Hash   // 输入: hash函数，默认 sha256
	Checksum       string      // 输出: checksum
}

func (h *HTTPReceiveFile) String() string {
	m := map[string]interface{}{
		"MaxReceiveSize": h.MaxReceiveSize,
		"DstDirMode":     h.DstDirMode,
		"DstFileMode":    h.MaxReceiveSize,
		"DstDir":         h.DstDir,
		"DstFilename":    h.DstFilename,
		"FieldName":      h.FieldName,
		"FileMIME":       h.FileMIME,
		"TotalSize":      h.TotalSize,
	}
	body, _ := json.Marshal(m) // nolint
	return string(body)
}

type ReceiveFileOption func(*HTTPReceiveFile)

func WithReceiveFileHashMD5() ReceiveFileOption {
	return func(uploadFile *HTTPReceiveFile) {
		uploadFile.Hash = md5.New()
	}
}
func WithReceiveFileHashSha1() ReceiveFileOption {
	return func(uploadFile *HTTPReceiveFile) {
		uploadFile.Hash = sha1.New()
	}
}
func WithReceiveFileHashSha512() ReceiveFileOption {
	return func(uploadFile *HTTPReceiveFile) {
		uploadFile.Hash = sha512.New()
	}
}

func WithReceiveFileHash(hash hash.Hash) ReceiveFileOption {
	return func(uploadFile *HTTPReceiveFile) {
		uploadFile.Hash = hash
	}
}

// WithReceiveFileMaxReceiveSize 最小单位是字节Bytes
func WithReceiveFileMaxReceiveSize(maxReceiveSize int64) ReceiveFileOption {
	return func(uploadFile *HTTPReceiveFile) {
		uploadFile.MaxReceiveSize = maxReceiveSize
	}
}

func WithReceiveFileDstDirMode(dstDirMode fs.FileMode) ReceiveFileOption {
	return func(uploadFile *HTTPReceiveFile) {
		uploadFile.DstDirMode = dstDirMode
	}
}

func WithReceiveFileDstFileMode(dstFileMode fs.FileMode) ReceiveFileOption {
	return func(uploadFile *HTTPReceiveFile) {
		uploadFile.DstFileMode = dstFileMode
	}
}

func WithReceiveFileDstFilename(dstFilename string) ReceiveFileOption {
	return func(uploadFile *HTTPReceiveFile) {
		uploadFile.DstFilename = dstFilename
	}
}

func WithReceiveFileDstFilenameFunc(fn func(*multipart.FileHeader) string) ReceiveFileOption {
	return func(uploadFile *HTTPReceiveFile) {
		uploadFile.DstFilename = fn(uploadFile.fileHeader)
	}
}

// WithReceiveFileDstDir 接收文件夹路径，默认是当前工作路径
func WithReceiveFileDstDir(dstDir string) ReceiveFileOption {
	return func(uploadFile *HTTPReceiveFile) {
		uploadFile.DstDir = dstDir
	}
}

// NewHttpReceiveFile 创建http文件接收函数
func NewHttpReceiveFile(options ...ReceiveFileOption) func(r *http.Request) (HTTPReceiveFile, error) {
	return func(r *http.Request) (HTTPReceiveFile, error) {
		h := HTTPReceiveFile{
			MaxReceiveSize: GB,
			FieldName:      "file",
			DstDirMode:     os.ModePerm,
			DstFileMode:    os.FileMode(0755),
			DstDir:         ".",
			DstFilename:    "", // 默认使用上传文件名
			FileMIME:       r.Header.Get("Content-Type"),
			Hash:           sha256.New(),
		}
		pwd, err := os.Getwd()
		if err != nil {
			return h, err
		}
		h.DstDir = pwd
		srcFile, info, err := r.FormFile("file")
		if err != nil {
			return h, err
		}
		defer srcFile.Close() //nolint
		h.fileHeader = info
		h.DstFilename = info.Filename
		for _, opt := range options {
			opt(&h)
		}
		size, err := getSize(srcFile)
		if err != nil {
			return h, err
		}
		if size > h.MaxReceiveSize {
			return h, fmt.Errorf("file size %d Bytes exceeds the maximum value of %d Bytes", size, h.MaxReceiveSize)
		}
		if !FileExist(h.DstDir) {
			err = os.MkdirAll(h.DstDir, h.DstFileMode)
			if err != nil {
				return h, err
			}
		}
		if h.DstFilename == "" {
			return h, errors.New("filename cannot be empty")
		}
		tempFile, err := ioutil.TempFile(h.DstDir, "upload_")
		if err != nil {
			return h, err
		}
		h.Hash.Reset()
		h.TotalSize, err = io.Copy(io.MultiWriter(tempFile, h.Hash), srcFile)
		if err != nil {
			return h, err
		}
		err = tempFile.Close()
		if err != nil {
			return h, err
		}
		h.Checksum = hex.EncodeToString(h.Hash.Sum(nil))
		h.Hash.Reset()
		err = os.Rename(tempFile.Name(), filepath.Join(h.DstDir, h.DstFilename))
		if err != nil {
			os.Remove(tempFile.Name()) // nolint
			return h, err
		}
		return h, nil
	}
}

func getSize(content io.Seeker) (int64, error) {
	size, err := content.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	_, err = content.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return size, nil
}
