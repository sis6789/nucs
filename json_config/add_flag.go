package json_config

import (
	"flag"
	"log"
	"strings"
)

// AddFlag - 소스 전체에 정의된 flag를 시스템 제어 요소로 추가한다.
// 중복된 것은 최근의 값을 사용한다.
func AddFlag() {
	if !flag.Parsed() {
		flag.Parse()
	}
	// 전체 flag 값을 추가한다.
	flag.VisitAll(func(f *flag.Flag) {
		name := f.Name
		if !strings.HasPrefix(name, "test.") {
			value := f.Value
			Put(name, value)
			log.Printf("config\tflag\t%v\t%v", nStr(name), value)
		}
	})
}
