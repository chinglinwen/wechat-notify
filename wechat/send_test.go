package wechat

import (
	"flag"
	"fmt"
	"testing"
)

func init() {
	flag.Parse()
}

var (
	testSecret = "3Kds9ib-5JwY7-DrlxGIBq7XOjYDf846W3_Tda2sLe0"
	testParty  = "10"
)

func TestGenUsersString(t *testing.T) {
	u, err := genUsersString(testSecret, testParty, "wen")
	if err != nil {
		t.Error("genUsersString err", err)
		return
	}
	fmt.Println("got users: ", u)
}
func TestListmember(t *testing.T) {
	u, err := listmember(testSecret, testParty)
	if err != nil {
		t.Error("listmember err", err)
		return
	}
	for _, v := range u.Userlist {
		fmt.Printf("id: %v, name: %v\n", v.Userid, v.Name)
	}
	// fmt.Printf("%#v\n", u)
}

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
	_, err := getToken("3Kds9ib-5JwY7-DrlxGIBq7XOjYDf846W3_Tda2sLe0")
	if err == nil {
		t.Errorf("get token expect err invalid credential,got nil error\n")
		return
	}

	token, err := getToken(*defaultSecret)
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
		Agentid: *defaultAgentID,
		Secret:  *defaultSecret,
	}
	out, err := send(b)
	if err != nil {
		t.Errorf("send err %v\n", err)
		return
	}
	fmt.Println("send reply", out) //1000002: 突然报警
	// Send("", "3", "hello3", "1000002", "") //1000002: 突然报警
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
