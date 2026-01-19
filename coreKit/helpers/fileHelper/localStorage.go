package fileHelper

import (
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stringx"
)

type LocalStorage struct {
	FileName   string                // 文件名
	SavePath   string                // 保存路径
	MaxSize    int64                 // 最大文件大小（MB）
	Files      multipart.File        // 文件流
	FileHeader *multipart.FileHeader // 文件头
}

type LocalOption func(l *LocalStorage)

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(savePath string, files multipart.File, fileHeader *multipart.FileHeader, opts ...LocalOption) *LocalStorage {
	l := &LocalStorage{
		SavePath:   savePath,
		Files:      files,
		FileHeader: fileHeader,
	}
	for _, opt := range opts {
		opt(l)
	}

	if l.FileName == "" {
		l.FileName = l.defaultFileName()
	}

	return l
}

// WithFileName 设置自定义文件名
func WithFileName(filename string) LocalOption {
	return func(l *LocalStorage) {
		l.FileName = filename
	}
}

// WithMaxSize 设置最大文件大小（单位：MB）
func WithMaxSize(maxSize int64) LocalOption {
	return func(l *LocalStorage) {
		if maxSize > 0 {
			l.MaxSize = maxSize
		}
	}
}

// defaultFileName 生成默认文件名
func (l *LocalStorage) defaultFileName() string {
	ext := path.Ext(l.FileHeader.Filename)
	fileName := stringx.Rand() + "-" + time.Now().Format("20060102150405") + ext
	return fileName
}

// CheckFileExts 检查文件扩展名是否在允许的列表中
func (l *LocalStorage) CheckFileExts(fileHeader *multipart.FileHeader, allowedExts []string) bool {
	ext := path.Ext(fileHeader.Filename)
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// Save 保存文件到本地存储
func (l *LocalStorage) Save() (string, error) {
	if l.MaxSize > 0 {
		if CheckMaxSize(l.Files, int(l.MaxSize*1024*1024)) {
			return "", errors.New("超过文件上传限制")
		}
	}

	// 检查目录是否存在
	if !IsExists(l.SavePath) {
		// 创建文件夹
		if err := CreateSavePath(l.SavePath, os.ModePerm); err != nil {
			return "", errors.Wrapf(err, "无法创建保存目录：%s", l.SavePath)
		}
	}

	// 检查是否有写入权限
	if !CheckPermission(l.SavePath) {
		return "", errors.New("写入权限不够")
	}

	dst := filepath.Join(l.SavePath, l.FileName)
	return dst, SaveFile(l.FileHeader, dst)
}
