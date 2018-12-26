package wechat

import (
	"fmt"
	"testing"
)

func TestSend(t *testing.T) {
	fmt.Println(Send("wenzhenglin", "hello1", "1000002")) //1000002: 突然报警
}

func TestSends(t *testing.T) {
	fmt.Println(Sends([]string{"wenzhenglin", "wenzhenglin"}, "hello1 from sends", "1000002"))
}
