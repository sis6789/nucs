package like

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/biogo/biogo/align"
	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/seq/linear"
)

var smith = align.SW{
	{0, -1, -1, -1, -1},
	{-1, 2, -1, -1, -1},
	{-1, -1, 2, -1, -1},
	{-1, -1, -1, 2, -1},
	{-1, -1, -1, -1, 2},
}
var splitExpStr = `\[(\d+),(\d+)\)/\[(\d+),(\d+)\)=(-?\d+)|\[(\d+),(\d+)\)/(-)=(-?\d+)|(-)/\[(\d+),(\d+)\)=(-?\d+)`
var splitExp = regexp.MustCompile(splitExpStr)

func Like(s1, s2 string) (
	sStart int,
	ratio float32,
	qStart int,
	checkLen int,
	countMatch int,
	countFault int,
	sMatch string,
	qMatch string) {
	// use smith waterman
	fsa := &linear.Seq{Seq: alphabet.BytesToLetters([]byte(s1))}
	fsa.Alpha = alphabet.DNAgapped
	fsb := &linear.Seq{Seq: alphabet.BytesToLetters([]byte(s2))}
	fsb.Alpha = alphabet.DNAgapped

	aln, err := smith.Align(fsa, fsb)
	if err != nil {
		return -1, 0.0, -1, -1, -1, -1, "", ""
	}

	status := fmt.Sprint(aln)
	token := splitExp.FindAllStringSubmatch(status, -1)
	checkLen = 0
	countFault = 0
	countMatch = 0
	sStart = -1
	qStart = -1
	for _, t := range token {
		var match, fault, tkLen = 0, 0, 0
		switch {
		case t[8] == "-":
			// insert
			vf, _ := strconv.Atoi(t[6])
			vt, _ := strconv.Atoi(t[7])
			tkLen = vt - vf
			tkLen = 0
			break
		case t[10] == "-":
			// delete
			vf, _ := strconv.Atoi(t[11])
			vt, _ := strconv.Atoi(t[12])
			if qStart == -1 {
				qStart = vf
			}
			tkLen = vt - vf
			fault = tkLen
			checkLen += tkLen
			countFault += tkLen
			break
		default:
			// match and substitution
			vf, _ := strconv.Atoi(t[1])
			vt, _ := strconv.Atoi(t[2])
			if sStart == -1 {
				sStart = vf
			}
			vQ, _ := strconv.Atoi(t[3])
			if qStart == -1 {
				qStart = vQ
			}
			vPoint, _ := strconv.Atoi(t[5])
			tkLen = vt - vf
			match = (vPoint + tkLen) / 3
			fault = tkLen - match
			checkLen += tkLen
			countFault += fault
			countMatch += match
		}
	}
	ratio = float32(countMatch) / float32(checkLen)
	fa := align.Format(fsa, fsb, aln, '-')
	sMatch = fmt.Sprint(fa[0])
	qMatch = fmt.Sprint(fa[1])
	return sStart, ratio, qStart, checkLen, countMatch, countFault, sMatch, qMatch
}

func LikeSW(s1, s2 []byte) (score, start1, end1, start2, end2 int) {

	type matched struct {
		s1    int
		e1    int
		s2    int
		e2    int
		score int
	}
	type scorer interface{ Score() int }

	sws1 := &linear.Seq{Seq: alphabet.BytesToLetters(s1)}
	sws1.Alpha = alphabet.DNAgapped
	sws2 := &linear.Seq{Seq: alphabet.BytesToLetters(s2)}
	sws2.Alpha = alphabet.DNAgapped

	alignResult, _ := smith.Align(sws1, sws2)
	highestMatch := matched{score: -1000}
	for _, v := range alignResult {
		m1 := matched{
			s1:    v.Features()[0].Start(),
			e1:    v.Features()[0].End(),
			s2:    v.Features()[1].Start(),
			e2:    v.Features()[1].End(),
			score: v.(scorer).Score(),
		}
		if m1.score > highestMatch.score {
			highestMatch = m1
		}
	}

	return highestMatch.score, highestMatch.s1, highestMatch.e1, highestMatch.s2, highestMatch.e2
}
