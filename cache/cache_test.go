package cache

import (
	"fmt"
	"testing"
)

func TestCache(t *testing.T) {
	user, toparty, content, status, expire := "u", "3,", "body", "ok", "1m"
	Set(user, toparty, content, status, expire)
	a, found := Get(user, toparty, content, status)
	fmt.Println("get", a, found)
	if !found {
		t.Errorf("cache should be exist")
		return
	}
	_, found = Get(user, toparty, content, "nook")
	if found {
		t.Errorf("cache should not exist")
	}
}
