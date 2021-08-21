package json_config

import (
	"encoding/json"
	"log"
)

func (x *JsonConfig) Decode(s []byte) {
	x.m = make(map[string]interface{})
	if err := json.Unmarshal(s, &x.m); err != nil {
		log.Panicln(err)
	}
}

func (x JsonConfig) Int64(field string) int64 {
	v, exist := x.m[field]
	if exist {
		fv := v.(float64)
		iv := int64(fv)
		return iv
	}
	return 0
}
func (x JsonConfig) Float64(field string) float64 {
	v, exist := x.m[field]
	if exist {
		fv := v.(float64)
		return fv
	}
	return 0.0
}
func (x JsonConfig) String(field string) string {
	v, exist := x.m[field]
	if exist {
		s := v.(string)
		return s
	}
	return ""
}
