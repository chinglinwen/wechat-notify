package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

var wd string

// Where get file location and line number info
func Where() string {
	_, f, l, _ := runtime.Caller(1)
	return fmt.Sprintf("%v:%v", strings.TrimPrefix(f, wd+"/"), l)
}

func init() {
	wd, _ = os.Getwd()
	if runtime.GOOS == "windows" {
		wd = strings.Replace(wd, "\\", "/", -1)
	}
}
