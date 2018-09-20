package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"wen/wechat-notify/wechat"

	"github.com/chinglinwen/log"
	"github.com/namsral/flag"
)

var (
	defaultExpire         = flag.String("e", "10m", "default expire time")
	defaultExpireDuration time.Duration
)

func init() {
	flag.Parse()

	var err error
	defaultExpireDuration, err = time.ParseDuration(*defaultExpire)
	if err != nil {
		log.Fatalf("default expire: %v parse err: %v\n", *defaultExpire, err)
	}
}

// test
// curl "localhost:8001/?user=wenzhenglin&content=aa"
func sendmsg(w http.ResponseWriter, req *http.Request) {
	user := req.FormValue("user")
	users := strings.Split(user, ",")
	expire := req.FormValue("expire")

	content := req.FormValue("content")
	if user == "" || content == "" {
		e := fmt.Sprintf("user: %v ,or content: %v is empty\n", user, content)
		log.Printf(e)
		fmt.Fprintf(w, e)
		return
	}

	if d, found := cacheGet(user, content); found {
		e := fmt.Sprintf("user: %v, content: %v not expired in %v, skip send\n",
			user, content, d.Format("15:04:05"))
		log.Printf(e)
		fmt.Fprintf(w, e)
		return
	}

	var d time.Duration
	if expire != "" {
		var err error
		d, err = time.ParseDuration(expire)
		if err != nil {
			e := fmt.Errorf("user: %v, expire: %v parse err: %v, ignore it\n", user, expire, err)
			log.Println(e)
		}
	} else {
		d = defaultExpireDuration
	}
	log.Printf("user %v expire set to %v\n", user, d)
	cacheSet(user, content, d)

	var msg string
	_, err := wechat.Sends(users, content)
	if err != nil {
		msg = fmt.Sprintf("send err: %v\n", err)
		log.Printf(msg)
		fmt.Fprintf(w, msg)
		return
	}
	msg = fmt.Sprintf("send ok: %v, %v\n", user, content)
	log.Printf(msg)
	fmt.Fprintf(w, msg)
	return
}

func main() {
	log.Println("starting...")
	http.HandleFunc("/", sendmsg)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Printf("start server error , err=%v", err)
	}
}
