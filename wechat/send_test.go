package wechat

import (
	"fmt"
	"testing"
)

func TestSend(t *testing.T) {
	fmt.Println(Send("wenzhenglin", "hello1"))
}

func TestSends(t *testing.T) {
	fmt.Println(Sends([]string{"wenzhenglin", "wenzhenglin"}, "hello1 from sends"))
}
