package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/SAIKAII/chatroom-backend/router"
)

func init() {
	dir, err := os.Stat("./files")
	if os.IsExist(err) {
		if dir.IsDir() {
			return
		}
	}

	if err := os.MkdirAll("./files", os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("Distribute Chat App v0.01")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
