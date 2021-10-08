package caller

import (
	"path/filepath"
	"runtime"
	"strconv"
)

func Caller() string {
	ix := 0
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
	return stack
}
