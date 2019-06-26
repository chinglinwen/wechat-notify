package wechat

import (
	"flag"
	"log"
)

var (
	debugFlag = flag.Bool("debug", false, "enable debug")

	defaultAgentID = flag.String("agentid", "", "agent id") //default agentid 告警机器人-运维
	defaultSecret  = flag.String("secret", "", "secret")
	defaultCorpID  = flag.String("corpid", "", "corp id")

	requestTokenHeader = flag.String("geturl", "https://qyapi.weixin.qq.com/cgi-bin/gettoken?", "token get url")
	pushHeader         = flag.String("accessurl", "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=", "token access url")
	simplelistHeader   = flag.String("memberlist", "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist?access_token=", "simplelist access url")
)

func CheckFlag() {
	if *defaultAgentID == "" {
		log.Fatal("agentid is empty")
	}
	if *defaultSecret == "" {
		log.Fatal("secret is empty")
	}
	if *defaultCorpID == "" {
		log.Fatal("corpid is empty")
	}

}
