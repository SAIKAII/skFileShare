package router

import (
	"net/http"

	"github.com/SAIKAII/chatroom-backend/app/api"
	"github.com/SAIKAII/chatroom-backend/pkg/websocket"
)

func init() {
	// 用于注册websocket聊天连接路由
	pool := websocket.NewPool()
	go pool.Start()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		api.ServeWs(pool, w, r)
	})

	http.HandleFunc("/getfile/", api.GetSpecifiedFile)
	http.HandleFunc("/getfiles/", api.GetFilesInfo)
	http.HandleFunc("/upload", api.UploadFile)
	http.HandleFunc("/delete/", api.RemoveFile)
	// mime.AddExtensionType(".js", "text/javascript")
	// mime.AddExtensionType(".css", "text/css")
	// http.Handle("/room/", http.StripPrefix("/room/", http.FileServer(http.Dir("template/"))))

	// 用于注册处理模糊匹配路由，如：获取指定文件api
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("static/"))))
}
