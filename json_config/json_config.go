package json_config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/sis6789/nucs/caller"
)

var jsonConfig = make(map[string]interface{})

func normalize() {
	// keep only lower case
	w := make(map[string]interface{})
	for k, v := range jsonConfig {
		k1 := strings.ToLower(k)
		k1 = strings.ReplaceAll(k1, "-", "")
		k1 = strings.ReplaceAll(k1, "_", "")
		w[k1] = v
	}
	jsonConfig = w
}

func nStr(w string) string {
	k1 := strings.ToLower(w)
	k1 = strings.ReplaceAll(k1, "-", "")
	k1 = strings.ReplaceAll(k1, "_", "")
	return k1
}

// Map - return address of map structure "map[string]interface{}"
func Map() *map[string]interface{} {
	return &jsonConfig
}

// Decode : 입력을 json으로 바꾸고 해당 K,v 를 저장하거나 치환한다.
func Decode(s []byte) {
	tM := make(map[string]interface{})
	if err := json.Unmarshal(s, &tM); err != nil {
		log.Printf("%v %v", caller.Caller(), err)
		return
	}
	Set(tM)
	normalize()
}

// Encode : 저장 자료 전체를 JSON string으로 반환한다.
func Encode() string {
	if bs, err := json.Marshal(jsonConfig); err == nil {
		return string(bs)
	} else {
		return ""
	}
}

// Read : JSON 파일을 읽어 환경을 설정한다.
func Read(fileName string) {
	var err error
	var fBytes []byte
	fBytes, err = ioutil.ReadFile(fileName)
	if err != nil {
		log.Panicln(err)
	}
	tM := make(map[string]interface{})
	if err := json.Unmarshal(fBytes, &tM); err != nil {
		log.Printf("%v %v", caller.Caller(), err)
		return
	}
	Set(tM)
	normalize()
}

// Write : 현재의 설정 상태를 파일에 저장한다.
func Write(fileName string) {
	var err error
	var fBytes []byte
	fBytes, err = json.Marshal(jsonConfig)
	if err != nil {
		log.Printf("%v %v", caller.Caller(), err)
	}
	err = ioutil.WriteFile(fileName, fBytes, 0777)
	if err != nil {
		log.Printf("%v %v", caller.Caller(), err)
	}
}

// Exist - 필드 존재 확인
func Exist(field string) bool {
	_, exist := jsonConfig[nStr(field)]
	return exist
}

// Get : 필드 값을 반환한다. 반환 값은 empty interface로 적절한 type inference를 해야 한다.
func Get(field string) interface{} {
	v, exist := jsonConfig[nStr(field)]
	if exist {
		return v
	} else {
		return nil
	}
}

// Put : 환경 값을 저장한다. 동일한 필드가 존재하면 값을 대체한다. 없으면 신설한다.
func Put(field string, v interface{}) {
	jsonConfig[nStr(field)] = v
	normalize()
}

// Set : 환경 값을 입력 값으로 치환한다.
func Set(setValue map[string]interface{}) {
	for k, v := range setValue {
		jsonConfig[k] = v
	}
	normalize()
}

// Int : 숫자 필드로서 Int64 값을 반환한다. 숫자 값이 아니면 panic한다.
func Int(field string) int {
	v, exist := jsonConfig[nStr(field)]
	if !exist {
		return 0
	}
	switch v.(type) {
	case int:
		return v.(int)
	case float64:
		return int(v.(float64))
	case string:
		if fv, err := strconv.ParseFloat(v.(string), 64); err != nil {
			return 0
		} else {
			return int(fv)
		}
	default:
		return 0.0
	}
}

// Float64 : 숫자 필드로서 Float64 값을 반환한다. 숫자 값이 아니면 panic한다.
func Float64(field string) float64 {
	v, exist := jsonConfig[nStr(field)]
	if !exist {
		return 0.0
	}
	switch v.(type) {
	case float64:
		return v.(float64)
	case int:
		return float64(v.(int))
	case string:
		if fv, err := strconv.ParseFloat(v.(string), 64); err != nil {
			return 0.0
		} else {
			return fv
		}
	default:
		return 0.0
	}
}

// String : 필드값을 문자열로 반환한다.
func String(field string) string {
	v := jsonConfig[nStr(field)]
	return fmt.Sprint(v)
}

// Report : 환경 변수별 값을 개별 라인으로 반환한다.
func Report() string {
	var kl []string
	for k := range jsonConfig {
		kl = append(kl, k)
	}
	sort.Strings(kl)
	s := ""
	for _, k := range kl {
		s += fmt.Sprintf("%v\t%v\n", k, jsonConfig[k])
	}
	return s[0 : len(s)-1]
}

// ReportSlice : 환경 변수별 값을 string slice로 반환한다.
func ReportSlice() []string {
	var kl []string
	for k := range jsonConfig {
		kl = append(kl, k)
	}
	sort.Strings(kl)
	var s []string
	for _, k := range kl {
		s = append(s, fmt.Sprintf("%v\t%v", k, jsonConfig[k]))
	}
	return s
}

// File - 지정한 폴더에 지정한 파일에 대해 *os.File을 반환한다. (mode: c, a, r, rw)
func File(folder, name, mode string) *os.File {
	var err error
	var f *os.File
	dir := folder
	if !filepath.IsAbs(folder) {
		dir = filepath.Join(jsonConfig["workDir"].(string), folder)
	}
	if err = os.MkdirAll(dir, 0777); err != nil {
		log.Fatalln(caller.Caller(), err)
	}
	fullPath := filepath.Join(dir, name)
	switch mode {
	case "c", "create":
		f, err = os.Create(fullPath)
		if err != nil {
			log.Fatalln(caller.Caller(), err)
		}
		return f
	case "a", "append":
		f, err = os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			log.Fatalln(caller.Caller(), err)
		}
		return f
	case "r", "raad":
		f, err = os.Open(fullPath)
		if err != nil {
			log.Fatalln(caller.Caller(), err)
		}
		return f
	case "rw", "readwrite":
		f, err = os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			log.Fatalln(caller.Caller(), err)
		}
		return f
	}
	return nil
}
