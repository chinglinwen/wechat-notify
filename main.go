package main

import (
	"fmt"
	"net/http"
	"strings"
	"wen/wechat-notify/wechat"

	"github.com/chinglinwen/checkup/cache"
	"github.com/chinglinwen/log"
	"github.com/namsral/flag"
)

// agentid 2 is 告警机器人-运维

// the user is email prefix, there's no group
// send to many people through one by one?
// test
// curl -s "localhost:8001/?user=wenzhenglin&content=test"
// curl -s "localhost:8001/?user=wenzhenglin|zhaixg&content=test"
func sendmsg(w http.ResponseWriter, req *http.Request) {
	user := req.FormValue("user")
	users := strings.Split(user, ",")
	expire := req.FormValue("expire")

	content := req.FormValue("content")
	agentid := req.FormValue("agentid")
	status := req.FormValue("status")

	if user == "" || content == "" {
		e := fmt.Sprintf("user: %v ,or content: %v is empty\n", user, content)
		log.Printf(e)
		fmt.Fprintf(w, e)
		return
	}

	if d, found := cache.Get(user, content, status); found {
		e := fmt.Sprintf("user: %v, content: %v not expired in %v, skip send\n",
			user, content, d.Format("15:04:05"))
		log.Printf(e)
		fmt.Fprintf(w, e)
		return
	}

	// default should have no cache, send everytime
	if expire != "" {
		log.Printf("user %v,%v, status: %v, expire set to %v\n", user, content, status, expire)
		cache.Set(user, content, status, expire)
	}

	content += " " + status
	var msg string
	_, err := wechat.Sends(users, content, agentid)
	if err != nil {
		msg = fmt.Sprintf("send user: %v, %v, status: %v, expire: %v, err: %v\n", user, content, status, expire, err)
		log.Printf(msg)
		fmt.Fprintf(w, msg)
		return
	}
	msg = fmt.Sprintf("send ok: %v, %v, status: %v, expire: %v\n", user, content, status, expire)
	log.Printf(msg)
	fmt.Fprintf(w, msg)
	return
}

func main() {
	flag.Parse()
	log.Println("starting...")
	http.HandleFunc("/", sendmsg)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Printf("start server error , err=%v", err)
	}
}
