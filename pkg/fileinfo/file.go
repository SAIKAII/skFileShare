package fileinfo

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Type int    `json:"type"`
	Name string `json:"name"`
	Ext  string `json:"ext"`
	Size int64  `json:"size"`
}

type FileNotExistError struct {
	filename string
}

func (e FileNotExistError) Error() string {
	return fmt.Sprintf("文件：%s不存在", e.filename)
}

// GetFiles 返回匹配dirExp正则的所有文件信息
func GetFilesInfo(dirExp string) ([]FileInfo, error) {
	dirExp += "/*"
	files, err := filepath.Glob(dirExp)
	if err != nil {
		return nil, err
	}

	filesInfo := []FileInfo{}
	for _, v := range files {
		f, err := os.Stat(v)
		if err != nil {
			continue
		}
		filesInfo = append(filesInfo, FileInfo{
			Type: 11,
			Name: filepath.Base(v),
			Ext:  filepath.Ext(v),
			Size: f.Size(),
		})
	}

	return filesInfo, nil
}

// GetSpecifiedFile 获取指定文件，可以是绝对路径，也可以是相对路径
func GetSpecifiedFileInfo(filename string) (*FileInfo, error) {
	file, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, FileNotExistError{filename: filename}
		}

		return nil, err
	}

	f := &FileInfo{
		Type: 11,
		Name: filepath.Base(filename),
		Ext:  filepath.Ext(filename),
		Size: file.Size(),
	}
	return f, nil
}

// DownloadFile 用于读取指定文件并发送到客户端
func DownloadFile(filepath string, w io.Writer) error {
	file, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return FileNotExistError{filename: filepath}
		}

		return err
	}

	_, err = io.Copy(w, file)
	if err != nil {
		return err
	}

	return nil
}

// SaveFile 保存文件到指定路径并返回文件信息
func SaveFile(filename string, data []byte) (*FileInfo, error) {
	if err := ioutil.WriteFile(filename, data, os.ModePerm); err != nil {
		return nil, err
	}

	f, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	fileInfo := &FileInfo{
		Type: 11,
		Name: filepath.Base(filename),
		Ext:  filepath.Ext(filename),
		Size: f.Size(),
	}
	return fileInfo, nil
}
