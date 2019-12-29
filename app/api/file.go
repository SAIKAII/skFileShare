package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/SAIKAII/chatroom-backend/pkg/fileinfo"
)

type retData struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var fileRegexp *regexp.Regexp
var dirRegexp *regexp.Regexp

// GetSpecifiedFiles 用于下载指定文件
func GetSpecifiedFile(w http.ResponseWriter, r *http.Request) {
	path, err := url.QueryUnescape(r.RequestURI)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ret := fileRegexp.FindStringSubmatch(path)
	if len(ret) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p := filepath.Join("./files", ret[len(ret)-1])
	err = fileinfo.DownloadFile(p, w)
	if e, ok := err.(fileinfo.FileNotExistError); ok {
		w.WriteHeader(http.StatusNotFound)
		log.Println(e.Error())
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("获取文件时发生内部错误：%s\n", err)
		return
	}
}

// GetFiles 用于把共享目录下的所有文件信息遍历返回前端
func GetFiles(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	uri, err := url.QueryUnescape(r.RequestURI)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("遍历目录请求的参数出错")
		return
	}

	dir := dirRegexp.FindStringSubmatch(uri)
	var path string
	if len(dir) > 1 {
		path = filepath.Join("./files", dir[len(dir)-1])
	} else {
		path = "./files"
	}

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
		ret, _ := json.Marshal(&retData{
			Status:  400,
			Message: "请求有误",
		})
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
			ret, _ := json.Marshal(&retData{
				Status:  500,
				Message: "保存文件时发生了错误",
			})
			w.Write(ret)
			log.Printf("保存文件时出错，原因是：%s", err)
			return
		}

		defer file.Close()
		if _, err := io.Copy(file, f); err != nil {
			ret, _ := json.Marshal(&retData{
				Status:  500,
				Message: "保存文件时发生了错误",
			})
			w.Write(ret)
			log.Printf("写入文件时出错，原因是：%s", err)
			return
		}

		ret, _ := json.Marshal(&retData{
			Status:  200,
			Message: "文件已保存",
		})
		w.Write(ret)
	} else {
		ret, _ := json.Marshal(&retData{
			Status:  400,
			Message: "请求有误",
		})
		w.Write(ret)
		log.Printf("解析请求参数时出错，原因是：%s", e)
	}
}

func init() {
	var err error
	fileRegexp, err = regexp.Compile("^/getfile/(.*)$")
	if err != nil {
		log.Fatal(err)
	}

	dirRegexp, err = regexp.Compile("^/getfiles/?(.*)$")
	if err != nil {
		log.Fatal(err)
	}
}
