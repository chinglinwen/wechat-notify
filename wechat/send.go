package wechat

import (
	"fmt"
	"log"
	"strings"

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
   "agentid": "1",
   "text": {
       "content": "%v"
   },
   "safe":"0"
}
`

type body struct {
	touser  []string
	content string
}

const (
	CropID             = "wxbc5446c29a6e633f"
	Secret             = "MJTShf7z4sSqUNce0ywG9hmA1ek34jP1tnuVINO-i59N4RTmZslFEElM03OUdhEo"
	requestTokenHeader = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?"
	pushHeader         = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token="
)

func getToken() (string, error) {
	requestTokenUrl := fmt.Sprintf("%vcorpid=%v&corpsecret=%v", requestTokenHeader, CropID, Secret)
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

	return fmt.Sprintf("%v%v", pushHeader, token), nil
}

func genBody(b body) string {
	users := strings.Join(b.touser, "|")
	return fmt.Sprintf(bodytemplate, users, b.content)
}

func Send(user, content string) (string, error) {
	return send(body{
		touser:  []string{user},
		content: content,
	})
}

func Sends(users []string, content string) (string, error) {
	return send(body{
		touser:  users,
		content: content,
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
