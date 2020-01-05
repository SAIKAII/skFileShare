package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/SAIKAII/chatroom-backend/pkg/fileinfo"
)

type retData struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func newRetData(status int, message string) []byte {
	ret, err := json.Marshal(&retData{
		Status:  status,
		Message: message,
	})
	if err != nil {
		log.Printf("error: %s", err)
		return nil
	}

	return ret
}

// GetSpecifiedFile 用于下载指定文件
func GetSpecifiedFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.FormValue("file"))
	filename := r.FormValue("file")
	p := filepath.Join("./files/", filename)
	err := fileinfo.DownloadFile(p, w)
	if err != nil {
		if e, ok := err.(fileinfo.FileNotExistError); ok {
			w.WriteHeader(http.StatusNotFound)
			log.Println(e.Error())
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("获取文件时发生内部错误：%s\n", err)
		return
	}
}

// GetFilesInfo 用于把共享目录下的所有文件信息遍历返回前端
func GetFilesInfo(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm()
	// dir := r.FormValue("dir")

	// var path string
	// if len(dir) > 0 {
	// 	path = filepath.Join("./files", dir)
	// } else {
	// 	path = "./files"
	// }
	path := "./files"
	files, err := fileinfo.GetFilesInfo(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("遍历目录时发生错误：%s", err)
		return
	}

	result, err := json.Marshal(files)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("转换JSON时出错：%s", err)
		return
	}

	_, err = io.WriteString(w, string(result))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("返回目录下内容信息的时候发生错误：%s", err)
	}
}

// UploadFile 客户端上传文件到服务器，共享给其他用户
func UploadFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32); err != nil {
		ret := newRetData(400, "请求有误")
		if ret == nil {
			return
		}

		w.Write(ret)
		log.Printf("上传文件时解析出现错误：%s", err)
		return
	}

	if f, h, e := r.FormFile("upload-file"); e == nil {
		defer func() { _ = f.Close() }()
		name := filepath.Base(h.Filename)
		pathSub := r.PostFormValue("path")
		fileSavePath := "./files/" + pathSub + "/" + name
		log.Printf("文件保存路径：%s", fileSavePath)
		file, err := os.Create(fileSavePath)
		if err != nil {
			ret := newRetData(500, "保存文件时发生了错误")
			if ret == nil {
				return
			}

			w.Write(ret)
			log.Printf("保存文件时出错，原因是：%s", err)
			return
		}

		defer file.Close()
		if _, err := io.Copy(file, f); err != nil {
			ret := newRetData(500, "保存文件时发生了错误")
			if ret == nil {
				return
			}

			w.Write(ret)
			log.Printf("写入文件时出错，原因是：%s", err)
			return
		}

		ret := newRetData(200, "文件已保存")
		if ret == nil {
			return
		}

		w.Write(ret)
	} else {
		ret := newRetData(400, "请求有误")
		if ret == nil {
			return
		}

		w.Write(ret)
		log.Printf("解析请求参数时出错，原因是：%s", e)
	}
}

// RemoveFile 删除指定文件
func RemoveFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	file := r.FormValue("file")
	p := filepath.Join("./files", file)
	err := fileinfo.DeleteFile(p)
	if err != nil {
		if _, ok := err.(fileinfo.FileNotExistError); ok {
			ret := newRetData(404, "指定文件不存在")
			if ret == nil {
				return
			}

			w.Write(ret)
			return
		}
		ret := newRetData(500, "删除文件时发生错误")
		if ret == nil {
			return
		}

		w.Write(ret)
		log.Printf("删除文件时发生错误，原因是： %s", err)
		return
	}

	ret := newRetData(200, "文件已成功删除")
	if ret == nil {
		return
	}

	w.Write(ret)
}
