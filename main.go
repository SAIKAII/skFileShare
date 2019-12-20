package main

import (
	"fmt"
	"net/http"

	_ "github.com/SAIKAII/chatroom-backend/router"
)

func main() {
	fmt.Println("Distribute Chat App v0.01")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
