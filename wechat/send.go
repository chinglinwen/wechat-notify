package wechat

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"

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

// var (
// 	defaultAgentID     string
// 	defaultSecret      string
// 	defaultCorpID      string
// 	requestTokenHeader string
// 	pushHeader         string
// 	debugFlag          bool
// )

// // init default values
// func Init(agentidx, secretx, corpidx, requestTokenHeaderx, pushHeaderx string, debug bool) {
// 	defaultAgentID = agentidx
// 	defaultSecret = secretx
// 	defaultCorpID = corpidx
// 	requestTokenHeader = requestTokenHeaderx
// 	pushHeader = pushHeaderx

// 	debugFlag = debug
// }

func getToken(secret string) (string, error) {
	requestTokenUrl := fmt.Sprintf("%vcorpid=%v&corpsecret=%v", *requestTokenHeader, *defaultCorpID, secret)
	resp, err := resty.SetDebug(*debugFlag).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		R().
		Get(requestTokenUrl)
	if err != nil {
		return "", fmt.Errorf("get token err: %v", err)
	}
	errcode := gjson.Get(resp.String(), "errcode").Int()
	if errcode != 0 {
		errmsg := gjson.Get(resp.String(), "errmsg").String()
		return "", fmt.Errorf("gettoken err: %v", strings.Split(errmsg, ",")[0])
	}
	token := gjson.Get(resp.String(), "access_token").String()
	if token == "" {
		return "", fmt.Errorf("got empty token")
	}
	return token, nil
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

func Send(user, toparty, content, agentid, secret string) {
	if agentid == "" {
		agentid = *defaultAgentID
	}
	if secret == "" {
		secret = *defaultSecret
	}
	bodyChan <- Body{
		Touser:  []string{user},
		Toparty: toparty,
		Content: content,
		Agentid: agentid,
		Secret:  secret,
	}
	return
}

func Sends(users []string, toparty, content, agentid, secret, exceptme string) (err error) {
	if agentid == "" {
		agentid = *defaultAgentID
	}
	if secret == "" {
		secret = *defaultSecret
	}
	if exceptme != "" {
		users, err = genUsersString(secret, toparty, exceptme)
		if err != nil {
			return err
		}
		toparty = ""
	}
	bodyChan <- Body{
		Touser:  users,
		Toparty: toparty,
		Content: content,
		Agentid: agentid,
		Secret:  secret,
	}
	return
}

var bodyChan = make(chan Body, 1000) // buffer item 1000 for the combine of burst

func combine() {
	log.Println("starting combine goroutine in the background")

	msgs := make(chan Body)
	go func(in chan Body) {
		msg := make(map[string]Body)
		for {
			select {
			case v := <-bodyChan:
				key := fmt.Sprintf("%v%v", v.Touser, v.Toparty)
				if _, ok := msg[key]; !ok {
					msg[key] = v
				} else {
					t := msg[key]
					t.Content += "\n=======\n" + v.Content
					msg[key] = t
				}
			case <-time.After(2000 * time.Millisecond): // combine between 2 seconds
				for k, v := range msg {
					msgs <- v
					delete(msg, k)
				}

			}
		}

	}(bodyChan)

	for v := range msgs {
		// time.Sleep(1 * time.Second)
		log.Printf("combined msg for user: %v, party: %v, for the busrt\n", v.Touser, v.Toparty)
		_, err := send(v)
		if err != nil {
			log.Printf("combine err: %v\n", err)
		}
	}
}

func init() {
	go combine()
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
	resp, err := resty.SetDebug(*debugFlag).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		R().
		SetBody(string(data)).
		Post(pushurl)

	if err != nil {
		return "", err
	}
	errcode := gjson.Get(resp.String(), "errcode").Int()
	if errcode != 0 {
		errmsg := gjson.Get(resp.String(), "errmsg").String()
		return "", fmt.Errorf("send err: %v", strings.Split(errmsg, ",")[0])
	}
	return resp.String(), nil
}
