package main

import (
	"fmt"
	"log"
	"net/http"
	"wen/wechat-notify/wechat"
)

// test
// curl "localhost:8001/?user=wenzhenglin&content=aa"
func sendmsg(w http.ResponseWriter, req *http.Request) {
	user := req.FormValue("user")
	content := req.FormValue("content")
	if user == "" || content == "" {
		log.Printf("user: %v ,or content: %v is empty\n", user, content)
		fmt.Fprintf(w, "user: %v ,or content: %v is empty\n", user, content)
		return
	}
	_, err := wechat.Send(user, content)
	if err != nil {
		log.Printf("send err: %v\n", err)
		fmt.Fprintf(w, "send err: %v\n", err)
		return
	}
	log.Printf("send ok: %v, %v\n", user, content)
	fmt.Fprintf(w, "send ok: %v, %v\n", user, content)
	return
}

func main() {
	http.HandleFunc("/", sendmsg)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Printf("start server error , err=%v", err)
	}
}
