package json_config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type JsonConfig struct {
	m map[string]interface{}
}

// New : 새로은 config 개체를 생성한다.
func New() JsonConfig {
	return JsonConfig{
		m: make(map[string]interface{}),
	}
}

// NewFile : 새로은 config 개체를 생성하고 해당 파일의 JSON으로 설정한다.
func NewFile(fileName string) JsonConfig {
	newConfig := JsonConfig{
		m: make(map[string]interface{}),
	}
	newConfig.Read(fileName)
	return newConfig
}

// Decode : 저장 자료 전체를 입력 JSON으로 교체한다.
func (x *JsonConfig) Decode(s []byte) {
	x.m = make(map[string]interface{})
	if err := json.Unmarshal(s, &x.m); err != nil {
		log.Panicln(err)
	}
}

// Read : JSON 파일을 읽어 환경을 설정한다.
func (x *JsonConfig) Read(fileName string) {
	var err error
	var fBytes []byte
	fBytes, err = ioutil.ReadFile(fileName)
	if err != nil {
		log.Panicln(err)
	}
	x.m = make(map[string]interface{})
	if err := json.Unmarshal(fBytes, &x.m); err != nil {
		log.Panicln(err)
	}
}

// Write : 현재의 설정 상태를 파일에 저장한다.
func (x JsonConfig) Write(fileName string) {
	var err error
	var fBytes []byte
	fBytes, err = json.Marshal(x.m)
	if err != nil {
		log.Panicln(err)
	}
	err = ioutil.WriteFile(fileName, fBytes, 0777)
	if err != nil {
		log.Panicln(err)
	}
}

// Get : 필드 값을 반환한다. 반환 값은 empty interface로 적절한 type inference를 해야 한다.
func (x JsonConfig) Get(field string) interface{} {
	v, exist := x.m[field]
	if exist {
		return v
	} else {
		return nil
	}
}

// Put : 환경 값을 저장한다. 동일한 필드가 존재하면 값을 대체한다. 없으면 신설한다.
func (x *JsonConfig) Put(field string, v interface{}) {
	x.m[field] = v
}

// Int64 : 숫자 필드로서 Int64 값을 반환한다. 숫자 값이 아니면 panic한다.
func (x JsonConfig) Int64(field string) int64 {
	v, exist := x.m[field]
	if exist {
		fv := v.(float64)
		iv := int64(fv)
		return iv
	}
	return 0
}

// Float64 : 숫자 필드로서 Float64 값을 반환한다. 숫자 값이 아니면 panic한다.
func (x JsonConfig) Float64(field string) float64 {
	v, exist := x.m[field]
	if exist {
		fv := v.(float64)
		return fv
	}
	return 0.0
}

// String : 필드값을 문자열로 반환한다.
func (x JsonConfig) String(field string) string {
	v := x.m[field]
	return fmt.Sprint(v)
}
