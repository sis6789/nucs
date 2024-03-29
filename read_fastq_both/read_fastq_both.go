package read_fastq_both

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type ReadPair struct {
	isReady   bool
	isPair    bool
	path      string
	err       error
	onlySide  int
	fileName  [2]string
	fileBytes [2][]byte
	reader    [2]*bytes.Reader
	//file       [2]*os.File
	scan       [2]*bufio.Scanner
	lineNumber [2]int
	recCount   [2]int
	Text       [2]string
	Phred      [2]string
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
	var err error
	switch len(name) {
	case 0:
		log.Fatalf("no file is specified")
	case 1:
		x.isPair = false
		x.path = path
		x.fileName[0] = name[0]
		x.fileBytes[0], err = os.ReadFile(filepath.Join(path, name[0]))
		if err != nil {
			log.Fatalf("%v", err)
		}
		x.reader[0] = bytes.NewReader(x.fileBytes[0])
		if filepath.Ext(name[0]) == ".gz" {
			gzReader, err := gzip.NewReader(x.reader[0])
			if err != nil {
				log.Fatalf("%v", err)
			}
			x.scan[0] = bufio.NewScanner(gzReader)
		} else {
			x.scan[0] = bufio.NewScanner(x.reader[0])
		}
	case 2:
		x.isPair = true
		x.path = path
		x.fileName[0] = name[0]
		x.fileName[1] = name[1]
		x.isReady = true
		// read 1
		x.fileBytes[0], err = os.ReadFile(filepath.Join(path, name[0]))
		if err != nil {
			log.Fatalf("%v", err)
		}
		x.reader[0] = bytes.NewReader(x.fileBytes[0])
		if filepath.Ext(name[0]) == ".gz" {
			gzReader, err := gzip.NewReader(x.reader[0])
			if err != nil {
				log.Fatalf("%v", err)
			}
			x.scan[0] = bufio.NewScanner(gzReader)
		} else {
			x.scan[0] = bufio.NewScanner(x.reader[0])
		}
		// read 2
		x.fileBytes[1], err = os.ReadFile(filepath.Join(path, name[1]))
		if err != nil {
			log.Fatalf("%v", err)
		}
		x.reader[1] = bytes.NewReader(x.fileBytes[1])
		if filepath.Ext(name[0]) == ".gz" {
			gzReader, err := gzip.NewReader(x.reader[1])
			if err != nil {
				log.Fatalf("%v", err)
			}
			x.scan[1] = bufio.NewScanner(gzReader)
		} else {
			x.scan[1] = bufio.NewScanner(x.reader[1])
		}
	default:
		log.Printf("too many parameters. path and max 2 name.")
		x.isReady = false
	}
}

// KeyomicsFastqFileNameRegex
// $0 full Name, $1 Sample, $2 SamplePart, $3 Hypen, $4 Repeat,
// $5 Underline, $6 R or Null, $7 Read(1,2), $8 file Extension
var KeyomicsFastqFileNameRegex = regexp.MustCompile(`(?i)^(\w+)(\d+)(-)(\d+)(_)(R?)(\d+)(\.(fastq|fq))$`)

func MakeNamePair(name string) string {
	nameToken := KeyomicsFastqFileNameRegex.FindStringSubmatch(name)
	if nameToken == nil {
		log.Panicln("invalid keyomics file name format")
		return ""
	}
	r1Name := nameToken[1] + nameToken[2] + nameToken[3] + nameToken[4] + nameToken[5] + nameToken[6] + "1" + nameToken[8]
	r2Name := nameToken[1] + nameToken[2] + nameToken[3] + nameToken[4] + nameToken[5] + nameToken[6] + "2" + nameToken[8]
	return r1Name + ":" + r2Name
}

// OpenPair - 경로와 R1/R2 파일을 준비한다. 지정은 R1 파일로 하고 R2 파일은 규칙에 따라 정해진 이름을 사용한다.
func (x *ReadPair) OpenPair(path string, r1Name string) {

	r1NameSplit := KeyomicsFastqFileNameRegex.FindStringSubmatch(r1Name)
	if r1NameSplit == nil {
		log.Panicln("r1 r1Name format invalid.")
		return
	}
	r2Name := r1NameSplit[1] + r1NameSplit[2] + r1NameSplit[3] + r1NameSplit[4] + r1NameSplit[5] + r1NameSplit[6] + "2" + r1NameSplit[8]
	x.Open(path, r1Name, r2Name)
}

// Close - 사용중인 파일을 닫고 구조체를 초기화 한다.
func (x *ReadPair) Close() {
	x.isReady = false
	x.isPair = false
	x.path = ""
	x.err = nil
	x.fileName = [2]string{"", ""}
	x.scan = [2]*bufio.Scanner{nil, nil}
	x.lineNumber = [2]int{0, 0}
	x.recCount = [2]int{0, 0}
	x.Text = [2]string{"", ""}
	x.reader = [2]*bytes.Reader{nil, nil}
	x.fileBytes = [2][]byte{nil, nil}
}

// Next - 다음 fastq sequence를 읽어 구조체 Text 필드에 저장한다.
// R2 서열은 앞 뒤 순서를 바꿔 저장한다.
// R1, R2 중 어느 것이든 먼저 EOF 상황이면 이후 false를 반환한다.
func (x *ReadPair) Next() bool {

	// read R1
readR1:
	for {
		if !x.scan[0].Scan() {
			if x.scan[0].Err() != nil {
				panic(x.scan[0].Err())
			}
			return false
		}
		x.lineNumber[0]++
		switch x.lineNumber[0] % 4 {
		case 0:
			// phred
			x.Phred[0] = x.scan[0].Text()
			break readR1
		case 2:
			// nucleotide sequence
			x.Text[0] = x.scan[0].Text()
			x.recCount[0]++
		}
	}
	// read R2
	if !x.isPair {
		return true
	}
readR2:
	for {
		if !x.scan[1].Scan() {
			if x.scan[1].Err() != nil {
				panic(x.scan[1].Err())
			}
			return false
		}
		x.lineNumber[1]++
		switch x.lineNumber[1] % 4 {
		case 0:
			// phred
			x.Phred[1] = x.scan[1].Text()
			ReverseString(&x.Phred[1])
			break readR2
		case 2:
			// nucleotide sequence
			x.Text[1] = x.scan[1].Text()
			ReverseComplementString(&x.Text[1])
			x.recCount[1]++
		}
	}
	return true
}

// Number - 현 상태의 record number를 반환한다.
// 한 레코드는 4개의 줄로 구성되며 각 set의 두번째 줄을 수록한다.
// 지정된 순서가 존재하지 않으면 panic으로 종결한다.
// 역순의 지정도 허용되지 않는다.
// 저장되는 서열은 역순의 상보합 서열이다.
func (x *ReadPair) Number() int {
	return x.recCount[0]
}

// AtRec - R1,R2 파일에서 지정된 FastQ rec number를 읽는다.
// 한 레코드는 4개의 줄로 구성되며 각 set의 두번째 줄을 수록한다.
// 지정된 순서가 존재하지 않으면 panic으로 종결한다.
// 역순의 지정도 허용되지 않는다.
// 저장되는 서열은 역순의 상보합 서열이다.
func (x *ReadPair) AtRec(recOrdinal int) bool {

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
					log.Printf("EOF before requested ordinal")
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
					log.Printf("%v", x.scan[1].Err())
				}
				return false
			}
			x.lineNumber[1]++
			if x.lineNumber[1]%4 == 2 {
				x.Text[1] = x.scan[1].Text()
				ReverseComplementString(&x.Text[1])
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
func (x *ReadPair) At2Rec(recOrdinal int) bool {

	if !x.isPair {
		log.Printf("single file. Read2 is illegal.")
		return false
	}
	if recOrdinal < x.recCount[1] {
		log.Printf("reverse ordinal is not supported.")
		return false
	}

	for x.recCount[1] < recOrdinal {
		// read R2
		for {
			if !x.scan[1].Scan() {
				if x.scan[1].Err() != nil {
					log.Printf("%v", x.scan[1].Err())
				}
				if recOrdinal > x.recCount[1] {
					log.Printf("EOF before request ordinal")
				}
				return false
			}
			x.lineNumber[1]++
			if x.lineNumber[1]%4 == 2 {
				x.Text[1] = x.scan[1].Text()
				ReverseComplementString(&x.Text[1])
				x.recCount[1]++
				break
			}
		}
	}
	return true
}
