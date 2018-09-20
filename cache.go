package main

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

var C = cache.New(10*time.Minute, 10*time.Minute)

func cacheSet(user, content string, d time.Duration) {
	C.Set(user, content, d)
}

func cacheGet(user, content string) (time.Time, bool) {
	cc, d, found := C.GetWithExpiration(user)
	if found {
		if ccstr, ok := cc.(string); ok && ccstr == content {
			return d, true
		}
	}
	return time.Time{}, false
}
