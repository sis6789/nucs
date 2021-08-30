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

	if _, exist := jsonConfig["root_dir"]; !exist {
		jsonConfig["root_dir"], _ = homedir.Dir()
	}
	if _, exist := jsonConfig["job_title"]; !exist {
		jsonConfig["job_title"] = time.Now().Format("20060102")
	}
	if _, exist := jsonConfig["run_name"]; !exist {
		jsonConfig["run_name"] = time.Now().Format("20060102-1504")
	}
	if _, exist := jsonConfig["fastq_query_terminator"]; exist {
		jsonConfig["fastq_query_terminator_length"] = len(jsonConfig["fastq_query_terminator"].(string))
	} else {
		jsonConfig["fastq_query_terminator_length"] = 0
	}
	jsonConfig["work_dir"] = filepath.Join(jsonConfig["root_dir"].(string),
		jsonConfig["job_title"].(string), jsonConfig["run_name"].(string))
	jsonConfig["log_dir"] = filepath.Join(jsonConfig["work_dir"].(string), "log")
	jsonConfig["save_dir"] = filepath.Join(jsonConfig["work_dir"].(string), "save")
	jsonConfig["temp_dir"] = filepath.Join(jsonConfig["work_dir"].(string), "temp")
}
