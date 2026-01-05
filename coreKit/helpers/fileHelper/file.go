package fileHelper

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

// IsExists 检查目录或者文件是否存在
// true 存在，false 不存在
func IsExists(path string) bool {
	_, err := os.Stat(path) // 获取文件的描述信息 FileInfo
	return err == nil || os.IsExist(err)
}

// Put 将数据存入文件中
func Put(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

// NameWithoutExtension 返回不带文件后缀的文件名称
func NameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

// CheckMaxSize 检查文件大小是否超出最大限制（kb 进行比较）
// true 超出，false 没有超出
func CheckMaxSize(f multipart.File, maxSize int) bool {
	content, _ := io.ReadAll(f)
	return len(content) >= maxSize
}

// CheckPermission 检查文件权限是否足够
// true 足够 false 不够
func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)
	return !os.IsPermission(err)
}

// CreateSavePath 创建在上传文件时所使用的保存目录
func CreateSavePath(dst string, perm os.FileMode) error {
	// 该方法将会以传入的 os.FileMode 权限位去递归创建所需的所有目录结构，
	// 若涉及的目录均已存在，则不会进行任何操作，直接返回 nil
	return os.MkdirAll(dst, perm)
}

// SaveFile 保存文件
// 参考于：gin.Context.SaveUploadedFile() 方法
func SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open() // 打开源地址的文件
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst) // 创建目标地址的文件
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
