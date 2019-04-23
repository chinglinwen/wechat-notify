package wechat

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/namsral/flag"
	"github.com/tidwall/gjson"
	resty "gopkg.in/resty.v1"
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
	{{with .TouserRaw }}"touser": "{{.}}",{{end}}
	{{with .Toparty }}"toparty": "{{.}}",{{end}}
   "msgtype": "text",
   "agentid": "{{.Agentid}}",
   "text": {
       "content": "{{.Content}}"
   },
   "safe":"0"
}
`

type Body struct {
	Touser    []string
	TouserRaw string
	Toparty   string
	Agentid   string
	Secret    string
	Content   string
}

var (
	AgentID = flag.String("agentid", "1000002", "agent id") //default agentid 告警机器人-运维
	Secret  = flag.String("secret", "0G22tGXTEgr4eFAX1jxbHSVoXeWtZ8DmCW4LQcEnXvM", "secret")
	// AgentID = flag.String("agentid", "1000003", "agent id")
	// Secret  = flag.String("secret", "G5h7CTEqkBw-Fe3luf2JM8UNNJAcYTpbXvpveY7M3lg", "secret")  //k8s

	CorpID = flag.String("corpid", "ww89720c104a10253f", "corp id")

	requestTokenHeader = flag.String("geturl", "https://qyapi.weixin.qq.com/cgi-bin/gettoken?", "token get url")
	pushHeader         = flag.String("accessurl", "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=", "token access url")
)

func getToken(secret string) (string, error) {
	requestTokenUrl := fmt.Sprintf("%vcorpid=%v&corpsecret=%v", *requestTokenHeader, *CorpID, secret)
	resp, err := resty.
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		R().
		Get(requestTokenUrl)
	if err != nil {
		return "", fmt.Errorf("get token err: %v", err)
	}
	return gjson.Get(resp.String(), "access_token").String(), nil
}

func getPushUrl(secret string) (string, error) {
	token, err := getToken(secret)
	if err != nil {
		log.Println("token err: ", token)
		return "", err
	}

	return fmt.Sprintf("%v%v", *pushHeader, token), nil
}

func genBody(b Body) (result string, err error) {
	b.TouserRaw = strings.Join(b.Touser, "|")
	t := template.Must(template.New("body").Parse(bodytemplate))
	var buf bytes.Buffer
	err = t.Execute(&buf, b)
	return buf.String(), err
}

func Send(user, toparty, content, agentid, secret string) (string, error) {
	if agentid == "" {
		agentid = *AgentID
	}
	if secret == "" {
		secret = *Secret
	}
	return send(Body{
		Touser:  []string{user},
		Toparty: toparty,
		Content: content,
		Agentid: agentid,
		Secret:  secret,
	})
}

func Sends(users []string, toparty, content, agentid, secret string) (string, error) {
	if agentid == "" {
		agentid = *AgentID
	}
	if secret == "" {
		secret = *Secret
	}
	return send(Body{
		Touser:  users,
		Toparty: toparty,
		Content: content,
		Agentid: agentid,
		Secret:  secret,
	})
}

func send(b Body) (string, error) {
	data, err := genBody(b)
	if err != nil {
		return "", err
	}
	pushurl, err := getPushUrl(b.Secret)
	if err != nil {
		return "", err
	}
	resp, err := resty. //SetDebug(true).
				SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
				R().
				SetBody(string(data)).
				Post(pushurl)
	if err != nil {
		return "", err
	}
	return resp.String(), nil
}
