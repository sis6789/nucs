package read_fastq_both

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type ReadPair struct {
	isReady    bool
	isPair     bool
	path       string
	err        error
	onlySide   int
	fileName   [2]string
	file       [2]*os.File
	scan       [2]*bufio.Scanner
	lineNumber [2]int32
	recCount   [2]int32
	Text       [2]string
}

// New - 패어 파일을 공통으로 읽도록 하는 구조체를 새로이 생성한다.
// 모든 제어 정보는 구조체에 저장하여 go-routine 수행에 독립적이다.
// 사용법은 New - Open - Next - Close이다.
// 처리 오류가 생긱면 panic한다.
func New() *ReadPair {
	var rp ReadPair
	return &rp
}

// Open - 경로와 최대 2개의 파일을 지정하여 동기화된 읽기를 준비한다.
func (x *ReadPair) Open(path string, name ...string) {
	switch len(name) {
	case 0:
		panic("path,name is required.")
	case 1:
		x.isPair = false
		x.path = path
		x.fileName[0] = name[0]
		if x.file[0], x.err = os.Open(filepath.Join(path, name[0])); x.err != nil {
			panic(x.err)
		}
		x.scan[0] = bufio.NewScanner(x.file[0])
	case 2:
		x.isPair = true
		x.path = path
		x.fileName[0] = name[0]
		x.fileName[1] = name[1]
		if x.file[0], x.err = os.Open(filepath.Join(path, name[0])); x.err != nil {
			panic(x.err)
		}
		x.scan[0] = bufio.NewScanner(x.file[0])
		if x.file[1], x.err = os.Open(filepath.Join(path, name[1])); x.err != nil {
			panic(x.err)
		}
		x.scan[1] = bufio.NewScanner(x.file[1])
	default:
		panic("too many parameters. path and max 2 name.")
	}
	x.isReady = true
}

// $0 full Name, $1 Sample, $2 SamplePart, $3 Hypen, $4 Repeat, $5 Underline, $6 R or Null, $7 Read(1,2), $8 file Extension
var r2NameRegex = regexp.MustCompile(`^(\w+)(\d+)(-)(\d+)(_)(R?)(\d+)(\.(fastq|fq))$`)

// OpenPair - 경로와 R1/R2 파일을 준비한다. 지정은 R1 파일로 하고 R2 파일은 규칙에 따라 정해진 이름을 사용한다.
func (x *ReadPair) OpenPair(path string, r1Name string) {

	r1NameSplit := r2NameRegex.FindStringSubmatch(r1Name)
	if r1NameSplit == nil {
		log.Panicln("r1 r1Name format invalid.")
		return
	}
	r2Name := r1NameSplit[1] + r1NameSplit[2] + r1NameSplit[3] + r1NameSplit[4] + r1NameSplit[5] + r1NameSplit[6] + "2" + r1NameSplit[8]
	x.Open(path, r1Name, r2Name)
}

// Close - 사용중인 파일을 닫고 구조체를 초기화 한다.
func (x *ReadPair) Close() {
	if x.isReady {
		if x.err = x.file[0].Close(); x.err != nil {
			panic(x.err)
		}
		if x.isPair {
			if x.err = x.file[1].Close(); x.err != nil {
				panic(x.err)
			}
		}
	}
	x.isReady = false
	x.isPair = false
	x.path = ""
	x.err = nil
	x.fileName = [2]string{"", ""}
	x.file = [2]*os.File{nil, nil}
	x.scan = [2]*bufio.Scanner{nil, nil}
	x.lineNumber = [2]int32{0, 0}
	x.recCount = [2]int32{0, 0}
	x.Text = [2]string{"", ""}
}

// Next - 다음 fastq sequence를 읽어 구조체 Text 필드에 저장한다.
// R2 서열은 앞 뒤 순서를 바꿔 저장한다.
// R1, R2 중 어느 것이든 먼저 EOF 상황이면 이후 false를 반환한다.
func (x *ReadPair) Next() bool {

	// read R1
	for {
		if !x.scan[0].Scan() {
			if x.scan[0].Err() != nil {
				panic(x.scan[0].Err())
			}
			return false
		}
		x.lineNumber[0]++
		if x.lineNumber[0]%4 == 2 {
			x.Text[0] = x.scan[0].Text()
			x.recCount[0]++
			break
		}
	}
	// read R2
	if !x.isPair {
		return true
	}
	for {
		if !x.scan[1].Scan() {
			if x.scan[1].Err() != nil {
				panic(x.scan[1].Err())
			}
			return false
		}
		x.lineNumber[1]++
		if x.lineNumber[1]%4 == 2 {
			x.Text[1] = x.scan[1].Text()
			reverseComplementString(&x.Text[1])
			x.recCount[1]++
			break
		}
	}
	return true
}

// Number - 현 상태의 record number를 반환한다.
// 한 레코드는 4개의 줄로 구성되며 각 set의 두번째 줄을 수록한다.
// 지정된 순서가 존재하지 않으면 panic으로 종결한다.
// 역순의 지정도 허용되지 않는다.
// 저장되는 서열은 역순의 상보합 서열이다.
func (x *ReadPair) Number() int32 {
	return x.recCount[0]
}

// AtRec - R1,R2 파일에서 지정된 FastQ rec number를 읽는다.
// 한 레코드는 4개의 줄로 구성되며 각 set의 두번째 줄을 수록한다.
// 지정된 순서가 존재하지 않으면 panic으로 종결한다.
// 역순의 지정도 허용되지 않는다.
// 저장되는 서열은 역순의 상보합 서열이다.
func (x *ReadPair) AtRec(recOrdinal int32) bool {

	if recOrdinal < x.recCount[0] {
		panic("reverse ordinal is not supported.")
	}

	for x.recCount[0] < recOrdinal {
		// read R1
		for {
			if !x.scan[0].Scan() {
				if x.scan[0].Err() != nil {
					panic(x.scan[0].Err())
				}
				if recOrdinal > x.recCount[0] {
					panic("EOF before request ordinal")
				}
				return false
			}
			x.lineNumber[0]++
			if x.lineNumber[0]%4 == 2 {
				x.Text[0] = x.scan[0].Text()
				x.recCount[0]++
				break
			}
		}
		// read R2
		if !x.isPair {
			continue
		}
		for {
			if !x.scan[1].Scan() {
				if x.scan[1].Err() != nil {
					panic(x.scan[1].Err())
				}
				return false
			}
			x.lineNumber[1]++
			if x.lineNumber[1]%4 == 2 {
				x.Text[1] = x.scan[1].Text()
				reverseComplementString(&x.Text[1])
				x.recCount[1]++
				break
			}
		}
	}
	return true
}

// At2Rec - R2 파일에서 지정된 FastQ rec number를 읽는다.
// 한 레코드는 4개의 줄로 구성되며 각 set의 두번째 줄을 수록한다.
// 지정된 순서가 존재하지 않으면 panic으로 종결한다.
// 역순의 지정도 허용되지 않는다.
// 저장되는 서열은 역순의 상보합 서열이다.
func (x *ReadPair) At2Rec(recOrdinal int32) bool {

	if !x.isPair {
		panic("single file. Read2 is illegal.")
	}
	if recOrdinal < x.recCount[1] {
		panic("reverse ordinal is not supported.")
	}

	for x.recCount[1] < recOrdinal {
		// read R2
		for {
			if !x.scan[1].Scan() {
				if x.scan[1].Err() != nil {
					panic(x.scan[1].Err())
				}
				if recOrdinal > x.recCount[1] {
					panic("EOF before request ordinal")
				}
				return false
			}
			x.lineNumber[1]++
			if x.lineNumber[1]%4 == 2 {
				x.Text[1] = x.scan[1].Text()
				reverseComplementString(&x.Text[1])
				x.recCount[1]++
				break
			}
		}
	}
	return true
}
