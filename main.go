package main

import (
	"fmt"
	"net/http"
	"strings"

	"wen/wechat-notify/cache"
	"wen/wechat-notify/wechat"

	"github.com/chinglinwen/log"
	"github.com/namsral/flag"
)

var (
	addr = flag.String("a", ":8001", "listening address")
)

// agentid 2 is 告警机器人-运维

// the user is email prefix, there's no group
// send to many people through one by one?
// test
// curl -s "localhost:8001/?user=wenzhenglin&content=test"
// curl -s "localhost:8001/?user=wenzhenglin|zhaixg&content=test"
// curl -s "localhost:8001/?user=WenZhengLin|LuRenJia&content=test2&agentid=1000002"
// curl -s "localhost:8001/?toparty=2&content=test2&agentid=1000002"
// curl -s "localhost:8001/?toparty=3&content=test4&agentid=1000002&expire=1m"
// curl -s "localhost:8001/?user=wenzhenglin&content=test5&agentid=1000002&expire=1m"
// curl -s "localhost:8001/?user=中文名&content=test4"  //not ok, only wechat id works
func sendmsg(w http.ResponseWriter, req *http.Request) {
	user := req.FormValue("user")
	users := strings.Split(user, ",")
	toparty := req.FormValue("toparty")
	expire := req.FormValue("expire")

	content := req.FormValue("content")
	agentid := req.FormValue("agentid")
	status := req.FormValue("status")

	if user == "" && toparty == "" {
		e := fmt.Sprintf("user: %v, or toparty: %v is empty\n", user, toparty)
		log.Printf(e)
		fmt.Fprintf(w, e)
		return
	}

	if content == "" {
		e := fmt.Sprintf("content: %v is empty\n", content)
		log.Printf(e)
		fmt.Fprintf(w, e)
		return
	}

	if d, found := cache.Get(user, toparty, content, status); found {
		e := fmt.Sprintf("user: %v, content: %v not expired in %v, skip send\n",
			user, content, d.Format("15:04:05"))
		log.Printf(e)
		fmt.Fprintf(w, e)
		return
	}

	// default should have no cache, send everytime
	if expire != "" {
		log.Printf("user %v,%v, status: %v, expire set to %v\n", user, content, status, expire)
		cache.Set(user, toparty, content, status, expire)
	}

	content += " " + status
	var msg string
	_, err := wechat.Sends(users, toparty, content, agentid)
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
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Printf("start server error , err=%v", err)
	}
}
