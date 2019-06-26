package wechat

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/tidwall/gjson"
	resty "gopkg.in/resty.v1"
)

func genUsersString(secret, party, exceptme string) (users []string, err error) {
	u, err := listmember(secret, party)
	if err != nil {
		err = fmt.Errorf("listmember err: %v", err)
		return
	}
	users = make([]string, 0)
	for _, v := range u.Userlist {
		if v.Userid == exceptme || v.Name == exceptme {
			continue
		}
		users = append(users, v.Userid)
	}
	// users = strings.Join(s, "|")
	return
}

func getListMemberUrl(secret, party string) (string, error) {
	token, err := getToken(secret)
	if err != nil {
		log.Println("token err: ", token)
		return "", err
	}

	return fmt.Sprintf("%v%v&department_id=%v&fetch_child=1", *simplelistHeader, token, party), nil
}

func listmember(secret, party string) (u *Userlist, err error) {
	listurl, err := getListMemberUrl(secret, party)
	if err != nil {
		return
	}
	resp, err := resty.SetDebug(*debugFlag).
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		R().
		// SetBody(string(data)).
		Post(listurl)

	if err != nil {
		return
	}
	errcode := gjson.Get(resp.String(), "errcode").Int()
	if errcode != 0 {
		errmsg := gjson.Get(resp.String(), "errmsg").String()
		err = fmt.Errorf("send err: %v", strings.Split(errmsg, ",")[0])
		return
	}

	u = &Userlist{}
	err = json.Unmarshal(resp.Body(), u)
	if err != nil {
		err = fmt.Errorf("unmarshal err: %v, body: %v", err, resp.String())
		return
	}
	return
}

type Userlist struct {
	Errcode  int    `json:"errcode"`
	Errmsg   string `json:"errmsg"`
	Userlist []struct {
		Userid     string `json:"userid"`
		Name       string `json:"name"`
		Department []int  `json:"department"`
	} `json:"userlist"`
}
