package read_fastq_both

import (
	"io/ioutil"
	"log"
	"reflect"
	"regexp"
	"strconv"

	"github.com/sis6789/nucs/caller"
)

var complementTable [128]byte

func init() {
	complementTable['A'] = 'T'
	complementTable['T'] = 'A'
	complementTable['C'] = 'G'
	complementTable['G'] = 'C'
	complementTable['a'] = 't'
	complementTable['t'] = 'a'
	complementTable['c'] = 'g'
	complementTable['g'] = 'c'
}

func MatchNamedField(namedPattern *regexp.Regexp, source string) (rVal struct {
	GName string
	GNum  int
	FNum  int
	RNum  int
	Ext   string
}) {
	nullZero := func(w string) int {
		if w == "" {
			return 0
		}
		if v, err := strconv.Atoi(w); err == nil {
			return v
		} else {
			return 0
		}
	}
	match := namedPattern.FindStringSubmatch(source)
	if match == nil {
		rVal.Ext = "ERROR!"
		return rVal
	}
	for i, name := range namedPattern.SubexpNames() {
		switch name {
		case "gname":
			rVal.GName = match[i]
		case "gnum":
			rVal.GNum = nullZero(match[i])
		case "fnum":
			rVal.FNum = nullZero(match[i])
		case "rnum":
			rVal.RNum = nullZero(match[i])
		case "ext":
			rVal.Ext = match[i]
		}
	}
	return rVal
}

// PairList : fastq 파일에 대한 키오믹스 명영규칙에 따라 존재하는 파일의 R1, R2 쌍이 값을 결정하여 그 목록을 반환한다.
// 각 값은 semi-colon으로 분리하여 사용하도록 한다.
// R1 파일만 존재하면 semi-colon 없이 하나의 값만을 반환한다.
func PairList(path string, fileNamePattern string) []string {
	regTarget := regexp.MustCompile(fileNamePattern)
	// get list
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalln(caller.Caller(), err)
	}
	var fnList []string
	var fnPair = ""
	var fnOne = false
	for _, fi := range files {
		if fi.IsDir() {
			continue
		}
		if !regTarget.MatchString(fi.Name()) {
			continue
		}
		fields := MatchNamedField(regTarget, fi.Name())
		switch fields.RNum {
		case 1:
			if fnOne {
				fnList = append(fnList, fnPair)
			}
			fnPair = fi.Name()
			fnOne = true
		case 2:
			if fnOne {
				fnPair += ";" + fi.Name()
				fnList = append(fnList, fnPair)
				fnOne = false
				fnPair = ""
			}
		}
	}

	return fnList
}

// string을 역순으로 배열한다. 입력을 직접 변경한다.
func reverseComplementString(s *string) {
	sb := []byte(*s)
	n := len(sb)
	swap := reflect.Swapper(sb)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
	for ix := 0; ix < len(sb); ix++ {
		sb[ix] = complementTable[sb[ix]]
	}
	*s = string(sb)
}
