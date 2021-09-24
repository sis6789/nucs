package caller

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func Caller() string {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		baseName := filepath.Base(file)
		return fmt.Sprintf("%s:%d", baseName, line)
	} else {
		return "unknown:-1"
	}
}
