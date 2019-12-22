package router

import (
	"net/http"

	"github.com/SAIKAII/chatroom-backend/app/api"
	"github.com/SAIKAII/chatroom-backend/pkg/websocket"
)

// type routerInfo struct {
// 	pattern string
// 	// f       func(w http.ResponseWriter, r *http.Request)
// 	f http.Handler
// }

// var routePath = []routerInfo{
// 	routerInfo{
// 		pattern: "^/getfile/(.*)$",
// 		f:       http.HandlerFunc(api.GetSpecifiedFile),
// 	},
// 	routerInfo{
// 		pattern: "^/getfiles/?(.*)$",
// 		f:       http.HandlerFunc(api.GetFiles),
// 	},
// 	routerInfo{
// 		pattern: "^/$",
// 		f:       http.StripPrefix("/", http.FileServer(http.Dir("template/"))),
// 	},
// }

// route 使用正则路由转发
// func route(w http.ResponseWriter, r *http.Request) {
// 	isFound := false
// 	for _, p := range routePath {
// 		reg, err := regexp.Compile(p.pattern)
// 		if err != nil {
// 			continue
// 		}
//
// 		if reg.MatchString(r.URL.Path) {
// 			isFound = true
// 			// p.f(w, r)
// 			p.f.ServeHTTP(w, r)
// 		}
// 	}
// 	if !isFound {
// 		fmt.Fprint(w, "404 Page Not Found!")
// 	}
// }

func init() {
	// 用于注册websocket聊天连接路由
	pool := websocket.NewPool()
	go pool.Start()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		api.ServeWs(pool, w, r)
	})

	http.HandleFunc("/getfile/", api.GetSpecifiedFile)
	http.HandleFunc("/getfiles/", api.GetFiles)
	http.HandleFunc("/upload", api.UploadFile)
	// mime.AddExtensionType(".js", "text/javascript")
	// mime.AddExtensionType(".css", "text/css")
	// http.Handle("/room/", http.StripPrefix("/room/", http.FileServer(http.Dir("template/"))))

	// 用于注册处理模糊匹配路由，如：获取指定文件api
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("static/"))))
}
