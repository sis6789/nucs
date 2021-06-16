package tail_match

import (
	"regexp"
	"sort"
)

func lessInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

type found struct {
	mStart int
	mLen   int
}

// MatchAny Return longest submatch
func MatchAny(target, query string, lead, matchMinimum int) (isMatch bool, matchStart int, matchLen int) {
	if matchMinimum < lead {
		matchMinimum = lead
	}
	var founds []found
	matchOne := regexp.MustCompile(query[0:lead])
	checkPos := matchOne.FindAllStringIndex(target, -1)
	if len(checkPos) == 0 {
		return false, -1, 0
	}
	// check partial matching
	lenTarget := len(target)
	lenQuery := len(query)
	for ix := len(checkPos) - 1; ix >= 0; ix-- {
		start := checkPos[ix]
		checkLen := lessInt(lenTarget-start[0], lenQuery)
		if checkLen < matchMinimum {
			continue
		}
		for partLen := checkLen; partLen > lead; partLen-- {
			if target[start[0]+partLen-1] != query[partLen-1] {
				continue
			}
			if target[start[0]:start[0]+partLen] == query[0:partLen] {
				founds = append(founds, found{start[0], partLen})
			}
		}
	}
	if len(founds) == 0 {
		return false, -1, 0
	} else {
		sort.Slice(founds, func(i, j int) bool {
			if founds[i].mLen > founds[j].mLen {
				return true
			} else if founds[i].mLen == founds[j].mLen {
				if founds[i].mStart > founds[j].mStart {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		})
		return true, founds[0].mStart, founds[0].mLen
	}
}

// MatchTail return just query submatch at tail of target
func MatchTail(target, query string, lead, matchMinimum int) (isMatch bool, matchStart int, matchLen int) {
	if matchMinimum < lead {
		matchMinimum = lead
	}
	matchOne := regexp.MustCompile(query[0:lead])
	checkPos := matchOne.FindAllStringIndex(target, -1)
	if len(checkPos) == 0 {
		return false, -1, 0
	}
	// check partial matching
	lenTarget := len(target)
	lenQuery := len(query)
	for ix := len(checkPos) - 1; ix >= 0; ix-- {
		start := checkPos[ix]
		if start[0] < lenTarget-lenQuery {
			// 끝에서부터 질의 길이보다 더 딸어져서 종료함
			break
		}
		checkLen := lessInt(lenTarget-start[0], lenQuery)
		if checkLen < matchMinimum {
			continue
		}
		if target[start[0]:] == query[:checkLen] {
			return true, start[0], checkLen
		}
	}
	return false, -1, 0
}
