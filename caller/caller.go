package caller

import (
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func Caller() string {
	ix := 2
	stack := ""
	for {
		_, file, line, ok := runtime.Caller(ix)
		if ok {
			stack += "<" + filepath.Base(file) + ":" + strconv.Itoa(line)
			ix++
		} else {
			break
		}
	}
	pos := strings.LastIndex(stack, "<")
	return stack[:pos] + "$"
}
