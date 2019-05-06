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

func TestSend(t *testing.T) {
	// fmt.Println(Send("wenzhenglin", "2", "hello2", "1000002")) //1000002: 突然报警
	Send("", "3", "hello3", "1000002", "") //1000002: 突然报警
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
