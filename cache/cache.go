// cache to avoid repeat send notice message
package cache

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

var C = cache.New(10*time.Minute, 10*time.Minute)

type msg struct {
	content string
	status  string
	d       time.Duration
	n       int
}

// how to detect it's new ( by parameter )
// have cache ( status change is a new )
// compare the old status and new status
func Set(user, toparty, content, status, expire string) {
	d, _ := time.ParseDuration(expire)
	C.Set(user+toparty+content, msg{content, status, d, 0}, d)
	// fmt.Println("set ", user, toparty)
}

func Get(user, toparty, content, status string) (time.Time, bool) {
	cc, d, found := C.GetWithExpiration(user + toparty + content)
	if found {
		if cmsg, ok := cc.(msg); ok && cmsg.content == content {
			if cmsg.status == "healthy" && cmsg.n == 0 {
				// ignore the first time ok status
				return time.Time{}, true
			}
			if cmsg.status == status {
				// cache exist and status not changed, return true
				return d, true
			} else {
				// if status changed, update the cache
				C.Set(user+toparty+content, msg{content, status, cmsg.d, cmsg.n + 1}, cmsg.d)
			}
		}
	} else {
		if status == "healthy" {
			return time.Time{}, true
		}
	}
	return time.Time{}, false
}
