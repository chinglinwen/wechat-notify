package wechat

import (
	"fmt"
	"log"
	"strings"

	"github.com/namsral/flag"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
)

/* doc
http://qydev.weixin.qq.com/wiki/index.php?title=%E6%B6%88%E6%81%AF%E7%B1%BB%E5%9E%8B%E5%8F%8A%E6%95%B0%E6%8D%AE%E6%A0%BC%E5%BC%8F


{
   "touser": "UserID1|UserID2|UserID3",
   "toparty": " PartyID1 | PartyID2 ",
   "totag": " TagID1 | TagID2 ",
   "msgtype": "text",
   "agentid": 1,
   "text": {
       "content": "Holiday Request For Pony(http://xxxxx)"
   },
   "safe":0
}

*/

var bodytemplate = `
{
   "touser": "%v",
   "msgtype": "text",
   "agentid": "%v",
   "text": {
       "content": "%v"
   },
   "safe":"0"
}
`

type body struct {
	touser  []string
	agentid string
	content string
}

var (
	AgentID            = flag.String("agentid", "2", "agent id") //default agentid 告警机器人-运维
	CorpID             = flag.String("corpid", "wxbc5446c29a6e633f", "corp id")
	Secret             = flag.String("secret", "MJTShf7z4sSqUNce0ywG9hmA1ek34jP1tnuVINO-i59N4RTmZslFEElM03OUdhEo", "secret")
	requestTokenHeader = flag.String("geturl", "https://qyapi.weixin.qq.com/cgi-bin/gettoken?", "token get url")
	pushHeader         = flag.String("accessurl", "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=", "token access url")
)

func getToken() (string, error) {
	requestTokenUrl := fmt.Sprintf("%vcorpid=%v&corpsecret=%v", *requestTokenHeader, *CorpID, *Secret)
	resp, err := resty.R().Get(requestTokenUrl)
	if err != nil {
		return "", fmt.Errorf("get token err: %v", err)
	}
	return gjson.Get(resp.String(), "access_token").String(), nil
}

func getPushUrl() (string, error) {
	token, err := getToken()
	if err != nil {
		log.Println("token err: ", token)
		return "", err
	}

	return fmt.Sprintf("%v%v", *pushHeader, token), nil
}

func genBody(b body) string {
	users := strings.Join(b.touser, "|")
	return fmt.Sprintf(bodytemplate, users, b.agentid, b.content)
}

func Send(user, content, agentid string) (string, error) {
	if agentid == "" {
		agentid = *AgentID
	}
	return send(body{
		touser:  []string{user},
		content: content,
		agentid: agentid,
	})
}

func Sends(users []string, content, agentid string) (string, error) {
	if agentid == "" {
		agentid = *AgentID
	}
	return send(body{
		touser:  users,
		content: content,
		agentid: agentid,
	})
}

func send(b body) (string, error) {
	data := genBody(b)
	pushurl, err := getPushUrl()
	if err != nil {
		return "", err
	}
	resp, err := resty. //SetDebug(true).
				R().
				SetBody(string(data)).
				Post(pushurl)
	if err != nil {
		return "", err
	}
	return resp.String(), nil
}
