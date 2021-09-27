package json_config

import (
	_ "embed"
	"github.com/mitchellh/go-homedir"
	"path/filepath"
	"runtime"
	"time"
)

//go:embed linux.json
var linuxJson []byte

//go:embed windows.json
var windowsJson []byte

// KeyomicsBasic : 키오믹스 기본 환경을 정의한다.
func KeyomicsBasic() {
	if runtime.GOOS == "windows" {
		Decode(windowsJson)
	} else {
		Decode(linuxJson)
	}

	if _, exist := jsonConfig["rootDir"]; !exist {
		jsonConfig["rootDir"], _ = homedir.Dir()
	}
	if _, exist := jsonConfig["jobTitle"]; !exist {
		jsonConfig["jobTitle"] = time.Now().Format("20060102")
	}
	if _, exist := jsonConfig["runName"]; !exist {
		jsonConfig["runName"] = time.Now().Format("20060102-1504")
	}
	if _, exist := jsonConfig["fastqQueryTerminator"]; exist {
		jsonConfig["fastqQueryTerminatorLength"] = len(jsonConfig["fastqQueryTerminator"].(string))
	} else {
		jsonConfig["fastqQueryTerminatorLength"] = 0
	}
	jsonConfig["workDir"] = filepath.Join(jsonConfig["rootDir"].(string),
		jsonConfig["jobTitle"].(string), jsonConfig["runName"].(string))
	jsonConfig["logDir"] = filepath.Join(jsonConfig["workDir"].(string), "log")
	jsonConfig["saveDir"] = filepath.Join(jsonConfig["workDir"].(string), "save")
	jsonConfig["tempDir"] = filepath.Join(jsonConfig["workDir"].(string), "temp")
}
