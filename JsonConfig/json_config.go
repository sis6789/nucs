package json_config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type JsonConfig struct {
	m map[string]interface{}
}

//New
func New() JsonConfig {
	return JsonConfig{
		m: make(map[string]interface{}),
	}
}

func (x *JsonConfig) Set(field string, v interface{}) {
	x.m[field] = v
}

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
