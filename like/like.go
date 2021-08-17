package like

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/biogo/biogo/align"
	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/seq/linear"

	"github.com/sis6789/nucs/utility"
)

var smith = align.SW{
	{0, -1, -1, -1, -1},
	{-1, 2, -1, -1, -1},
	{-1, -1, 2, -1, -1},
	{-1, -1, -1, 2, -1},
	{-1, -1, -1, -1, 2},
}
var splitExpStr = `\[(\d+),(\d+)\)/\[(\d+),(\d+)\)=(-?\d+)|(-)/\[(\d+),(\d+)\)=(-?\d+)`
var splitExp = regexp.MustCompile(splitExpStr)

func Like(s1, s2 string) (start int, ratio float32, matchLen int, m1Str, m2Str string) {
	// use smith waterman
	fsa := &linear.Seq{Seq: alphabet.BytesToLetters([]byte(s1))}
	fsa.Alpha = alphabet.DNAgapped
	fsb := &linear.Seq{Seq: alphabet.BytesToLetters([]byte(s2))}
	fsb.Alpha = alphabet.DNAgapped

	aln, err := smith.Align(fsa, fsb)
	if err != nil {
		return -1, 0.0, 0, "", ""
	}

	status := fmt.Sprint(aln)
	token := splitExp.FindAllStringSubmatch(status, -1)
	var mCount = 0
	var mStart = utility.MaxIntValue
	var mDel = 0
	for tkPos, tk := range token {
		if tk[6] == "-" {
			// 삭제
			sStart, _ := strconv.ParseInt(tk[7], 10, 32)
			sEnd, _ := strconv.ParseInt(tk[8], 10, 32)
			mDel += int(sStart - sEnd)
			if tkPos < len(token)-1 {
				matchLen += int(sStart - sEnd)
			}
		} else {
			// 정합 또는 치환
			sStart, _ := strconv.ParseInt(tk[1], 10, 32)
			mStart = utility.MinInt(mStart, int(sStart))
			sEnd, _ := strconv.ParseInt(tk[2], 10, 32)
			sPoint, _ := strconv.ParseInt(tk[5], 10, 32)
			sLen := sEnd - sStart
			mCount += int(sLen+sPoint) / 2
			matchLen += int(sLen)
		}
	}
	fa := align.Format(fsa, fsb, aln, '-')
	m1Str = fmt.Sprint(fa[0])
	m2Str = fmt.Sprint(fa[1])
	ratio = float32(mCount) / float32(len(s2))
	start = mStart
	// 끝부분이 삭제 인지를 확인하여 전방 정합을 점검
	if token[len(token)-1][6] == "-" {
		ratio = float32(mCount) / float32(matchLen)
	}
	return start, ratio, matchLen, m1Str, m2Str
}
