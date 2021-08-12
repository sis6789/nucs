package applyBtop

import (
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var btopParser = regexp.MustCompile(`(\d+)?([ACGT-]{2})?`)

func reverseAny(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

// Complement 대칭되는 원소로 치환한다.
func Complement(slice interface{}) {
	switch slice.(type) {
	case *byte:
		values := slice.(*byte)
		switch *values {
		case 'A':
			*values = 'T'
		case 'C':
			*values = 'G'
		case 'G':
			*values = 'C'
		case 'T':
			*values = 'A'
		case 'a':
			*values = 't'
		case 'c':
			*values = 'g'
		case 'g':
			*values = 'c'
		case 't':
			*values = 'a'
		}
	case []byte:
		values := slice.([]byte)
		for ix := 0; ix < len(values); ix++ {
			Complement(&values[ix])
		}
	case []string:
		values := slice.([]string)
		for ix := 0; ix < len(values); ix++ {
			values[ix] = ComplementString(values[ix])
		}
	case [][]string:
		values := slice.([][]string)
		for ix := 0; ix < len(values); ix++ {
			Complement(values[ix])
		}
	default:
		log.Println("Invalid type complement")
	}
}

func ComplementString(inStr string) (outStr string) {
	inBytes := []byte(inStr)
	Complement(inBytes)
	return string(inBytes)
}

// 문자열의 앞과 뒤를 바꾸고 대칭되는 것으로 바꾼다.
func reverseComplementString(inStr string) (outStr string) {
	inBytes := []byte(inStr)
	reverseAny(inBytes)
	Complement(inBytes)
	return string(inBytes)
}

// ReverseString 문자열의 앞과 뒤를 바꾼다.
func ReverseString(inStr string) (outStr string) {
	inBytes := []byte(inStr)
	reverseAny(inBytes)
	return string(inBytes)
}

type BtopApplyResult struct {
	BpAddress     int
	QueryBP       []byte
	GenomeAddress []int
	QueryAddress  []int
	ModifyAddress [][]string
	Line1         []string
	Line2         []string
	RStart        int
	RLen          int
}

type BtopApplyRequest struct {
	Query  string
	Qstart int
	Sstart int
	Length int
	Btop   string
}

// ApplyBtop
// 지정된 범위 이외는 소문자로, 역방향 blast 결과는 부분 전도 및 상대치환.
// 전제 2개의 줄로 표현한다.
// 1번 줄에는 전체 검색을 표한한다.
// 부분검색의 경우 찾아지지 않은 부분은 소문자로 표시한다.
// 누락된 부분은 1번 줄에 '-'로 표시된다.
// 2번 줄에는 치환, 삽입, 누락을 표현한다.
// 치환은 원래 게놈 값을 표시한다.
// 누락은 원래 게놈값을 표시하고 1번줄에 `-`를 삽입한다.
// 삽입은 직전 위치에 삽입된 값을 소문자로 표시하고 원래 자리에서는 제거한다.
func ApplyBtop(btopRequest BtopApplyRequest) BtopApplyResult {
	query := btopRequest.Query
	qstart := btopRequest.Qstart
	sstart := btopRequest.Sstart
	length := btopRequest.Length
	btop := btopRequest.Btop

	var queryBP []byte
	var genomeAddress []int
	var queryAddress []int
	var modifyAddress = make([][]string, 0)
	var line1 []string
	var line2 []string
	var rStart int
	var rLen int

	wQPos := 1
	wSDelta := 1
	wSPos := sstart
	if length < 0 {
		wSDelta = -1 // 게놈 주소를 감소하는 방향으로 저장하도록 함
	}

	btopToken := btopParser.FindAllStringSubmatch(btop, -1)
	prefix := query[0 : qstart-1]
	ixNext := qstart - 1
	middle1 := ""
	middle2 := ""
	for _, btCell := range btopToken {
		if btCell[1] != "" {
			// 숫자 부분 - 게놈과 동일한 원소가 식별됨
			sameLen, _ := strconv.Atoi(btCell[1])
			for _, v := range []byte(query[ixNext : ixNext+sameLen]) {
				queryBP = append(queryBP, v)
				genomeAddress = append(genomeAddress, wSPos)
				wSPos += wSDelta
				queryAddress = append(queryAddress, wQPos)
				wQPos++
				modifyAddress = append(modifyAddress, []string{})
			}
			middle1 += query[ixNext : ixNext+sameLen]
			middle2 += strings.Repeat(".", sameLen)
			ixNext += sameLen
		}
		if btCell[2] != "" {
			// 변경 발생
			switch {
			case btCell[2][0] == '-':
				// 누락 - 게놈 위치의 원소가 누락된 것으로 판별, 게놈 주소를 차지하도록 하고 "-"로 누락을 표시함.
				queryBP = append(queryBP, '-')
				genomeAddress = append(genomeAddress, wSPos)
				wSPos += wSDelta
				queryAddress = append(queryAddress, wQPos-1)
				modifyAddress = append(modifyAddress, []string{btCell[2]})
				middle1 += "-"
				middle2 += btCell[2][1:]
			case btCell[2][1] == '-':
				// 삽입. 2번줄의 직전 자리에 소문자로 삽입을 표시.
				wQPos++
				modifyAddress = append(modifyAddress[:len(modifyAddress)-1], append(modifyAddress[len(modifyAddress)-1], btCell[2]))
				if strings.HasSuffix(middle2, ".") {
					middle2 = middle2[:len(middle2)-1] + strings.ToLower(btCell[2][0:1])
				} else {
					// 직전 또한 변화가 있음. 이 경우 합쳐서 *로 표시하여 주의만 환기시킴
					middle2 = middle2[:len(middle2)-1] + "*"
				}
				// 삽입된 자리는 문자는 무시하고 직후부터 처리하도록 한다.
				ixNext++
			default:
				// 치환
				queryBP = append(queryBP, query[ixNext])
				genomeAddress = append(genomeAddress, wSPos)
				wSPos += wSDelta
				queryAddress = append(queryAddress, wQPos)
				wQPos++
				modifyAddress = append(modifyAddress, []string{btCell[2]})
				middle1 += query[ixNext : ixNext+1]
				middle2 += btCell[2][1:2]
				ixNext++
			}
		}
	}
	suffix := query[ixNext:]
	if length < 0 {
		// complement
		Complement(queryBP)
		Complement(modifyAddress)
		// reverse display string
		middle1 = reverseComplementString(middle1)
		middle2 = reverseComplementString(middle2)
		line1 = []string{strings.ToLower(reverseComplementString(suffix)),
			middle1,
			strings.ToLower(reverseComplementString(prefix))}
		line2 = []string{strings.Repeat(".", len(suffix)),
			middle2,
			strings.Repeat(".", len(prefix))}
		rStart = len(suffix) + 1
		rLen = len(line1)
	} else {
		line1 = []string{strings.ToLower(prefix),
			middle1,
			strings.ToLower(suffix)}
		line2 = []string{strings.Repeat(".", len(prefix)),
			middle2,
			strings.Repeat(".", len(suffix))}
		rStart = len(prefix) + 1
		rLen = len(line1)
	}
	return BtopApplyResult{0, queryBP, genomeAddress,
		queryAddress, modifyAddress,
		line1, line2, rStart, rLen}
}
