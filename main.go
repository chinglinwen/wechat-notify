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

// https://work.weixin.qq.com/wework_admin/frame#apps/modApiApp/5629500139363788
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
//
// curl -s "localhost:8001/?toparty=2&content=test2&agentid=1000002"
// curl -s "localhost:8001/?toparty=2&content=test5"
// curl -s "localhost:8001/?toparty=2&content=test6&agentid=1000003&secret=G5h7CTEqkBw-Fe3luf2JM8UNNJAcYTpbXvpveY7M3lg"
func sendmsg(w http.ResponseWriter, req *http.Request) {
	user := req.FormValue("user")
	users := strings.Split(user, ",")
	toparty := req.FormValue("toparty")
	expire := req.FormValue("expire")

	precontent := req.FormValue("precontent")
	content := req.FormValue("content")
	agentid := req.FormValue("agentid")
	secret := req.FormValue("secret")

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
		e := fmt.Sprintf("skip send(not expired in %v), user: %v, content: %v\n",
			d.Format("15:04:05"), user, content)
		log.Printf(e)
		fmt.Fprintf(w, e)
		return
	}

	// default should have no cache, send everytime
	if expire != "" {
		log.Printf("user %v,%v, status: %v, expire set to %v\n", user, content, status, expire)
		cache.Set(user, toparty, content, status, expire)
	}

	contentbody := precontent + content + " " + status
	var msg string
	_, err := wechat.Sends(users, toparty, contentbody, agentid, secret)
	if err != nil {
		msg = fmt.Sprintf("send err: %v, user: %v, %v, status: %v, expire: %v\n", err, user, content, status, expire)
		log.Printf(msg)
		fmt.Fprintf(w, msg)
		return
	}
	msg = fmt.Sprintf("send ok, user: %v, %v, status: %v, expire: %v\n", user, content, status, expire)
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
