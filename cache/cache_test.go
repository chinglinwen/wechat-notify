package cache

import (
	"fmt"
	"testing"
)

func TestCache(t *testing.T) {
	content := `2019/04/04 11:23:06 send ok: wenzhenglin, 时间: 2019-4-4 11:21:32
	类别: Warning
	名字: default/172.31.90.51
	-----
	来源: Node
	原因: ImageGCFailed
	内容: (combined from similar events): wanted to free 2300926361 bytes, but freed 0 bytes space with errors in image deletion: [rpc error: code = Unknown desc = Error response from daemon: conflict: unable to delete f4cb5e83f0a4 (cannot be forced) - image is being used by running container 8b894819fa43, rp... (omited)`
	user, toparty, status, expire := "wenzhenglin", "", "ok", "1m"
	Set(user, toparty, content, status, expire)
	a, found := Get(user, toparty, content, status)
	fmt.Println("get", a, found)
	if !found {
		t.Errorf("cache should be exist")
		return
	}
	t.Log("a", a)
	_, found = Get(user, toparty, content, "nook")
	if found {
		t.Errorf("cache should not exist")
		return
	}
	t.Log("ok")
}
