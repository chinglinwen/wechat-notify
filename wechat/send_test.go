package wechat

import (
	"fmt"
	"testing"
)

func TestGenBody(t *testing.T) {
	a := Body{
		Touser:  []string{"a"},
		Toparty: "b",
		Content: "content",
		Agentid: "1",
	}
	fmt.Println(genBody(a)) //1000002: 突然报警
}

func TestGetToken(t *testing.T) {
	_, err := getToken("aa")
	if err == nil {
		t.Errorf("get token expect err invalid credential,got nil error\n")
		return
	}

	token, err := getToken(secret)
	if err != nil {
		t.Errorf("get token err: %v\n", err)
		return
	}
	fmt.Println("got token", token)
}

func TestSend(t *testing.T) {
	b := Body{
		Touser:  []string{"wenzhenglin"},
		Toparty: "",
		Content: "hello2",
		Agentid: agentID,
		Secret:  secret,
	}
	out, err := send(b)
	if err != nil {
		t.Errorf("send err %v\n", err)
		return
	}
	fmt.Println("send reply", out) //1000002: 突然报警
	// Send("", "3", "hello3", "1000002", "") //1000002: 突然报警
}

func TestSends(t *testing.T) {
	Sends([]string{"wenzhenglin", "wenzhenglin"}, "", "hello1 from sends", "1000002", "")
}

/*
test combine

while read l; do
  curl -s "localhost:8001/?user=wenzhenglin&content=test$l&agentid=1000002&expire=1m"
done<<eof
a
b
c
d
e
eof

*/
