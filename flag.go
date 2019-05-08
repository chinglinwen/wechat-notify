package main

import (
	"log"
	"wen/wechat-notify/wechat"

	"github.com/namsral/flag"
)

var (
	addr  = flag.String("a", ":8001", "listening address")
	debug = flag.Bool("debug", false, "enable debug")

	AgentID = flag.String("agentid", "", "agent id") //default agentid 告警机器人-运维
	Secret  = flag.String("secret", "", "secret")
	CorpID  = flag.String("corpid", "", "corp id")

	RequestTokenHeader = flag.String("geturl", "https://qyapi.weixin.qq.com/cgi-bin/gettoken?", "token get url")
	PushHeader         = flag.String("accessurl", "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=", "token access url")
)

func init() {
	flag.Parse()
	if *AgentID == "" {
		log.Fatal("agentid is empty")
	}
	if *Secret == "" {
		log.Fatal("secret is empty")
	}
	if *CorpID == "" {
		log.Fatal("corpid is empty")
	}

	wechat.Init(*AgentID, *Secret, *CorpID, *RequestTokenHeader, *PushHeader, *debug)
}
