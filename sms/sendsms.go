package sms

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/go-resty/resty"
)

const (
	//appurl     = "http://smsapi.qianbao.com/api/sms/send"
	appurl     = "http://192.168.1.238/api/sms/send"
	appName    = "yw"
	secure_key = "d8V>0_4K$c45khZ{FcqS'"
)

// currently not working
func send(phone, content string) {
	md5value := genmd5(phone)

	d := map[string]string{
		"phone":   phone,
		"content": content,
		"appName": appName,
	}
	//fmt.Println(d)
	data, _ := json.Marshal(d)

	resp, err := resty.R().
		SetHeader("Content-Type", "application/text; charset=utf-8").
		SetHeader("code", md5value).
		SetBody(url.QueryEscape(string(data))).
		Post(appurl)

	fmt.Println(resp, err)
}

func genmd5(phone string) string {
	h := md5.New()
	io.WriteString(h, appName+phone+secure_key)
	return string(h.Sum(nil))
}
