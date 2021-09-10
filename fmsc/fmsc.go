package fmsc

import "sort"

type FMSC struct {
	File      string `bson:"f"`
	Molecular string `bson:"m"`
	Side      string `bson:"s"`
	Count     int    `bson:"c"`
	Ordinal   []int  `bson:"o"`
}

// NewFMSC 새 FMSC 구조를 반환한다.
func NewFMSC(file string, molecular string, side string, count int, ordinal int) FMSC {
	return FMSC{
		File:      file,
		Molecular: molecular,
		Side:      side,
		Count:     count,
		R12UniSeq: []int{},
		Ordinal:   []int{ordinal},
	}
}

// SortFMS sort slice by File-Molecular-Side
func SortFMS(sliceFMSC []FMSC) {
	sort.Slice(sliceFMSC, func(i, j int) bool {
		switch {
		case sliceFMSC[i].File < sliceFMSC[j].File:
			return true
		case sliceFMSC[i].File > sliceFMSC[j].File:
			return false
		case sliceFMSC[i].Molecular < sliceFMSC[j].Molecular:
			return true
		case sliceFMSC[i].Molecular > sliceFMSC[j].Molecular:
			return false
		case sliceFMSC[i].Side < sliceFMSC[j].Side:
			return true
		case sliceFMSC[i].Side > sliceFMSC[j].Side:
			return false
		default:
			return false
		}
	})
}

// SortFMSC sort slice by File-Molecular-Side-Count
func SortFMSC(sliceFMSC []FMSC) {
	sort.Slice(sliceFMSC, func(i, j int) bool {
		switch {
		case sliceFMSC[i].File < sliceFMSC[j].File:
			return true
		case sliceFMSC[i].File > sliceFMSC[j].File:
			return false
		case sliceFMSC[i].Molecular < sliceFMSC[j].Molecular:
			return true
		case sliceFMSC[i].Molecular > sliceFMSC[j].Molecular:
			return false
		case sliceFMSC[i].Side < sliceFMSC[j].Side:
			return true
		case sliceFMSC[i].Side > sliceFMSC[j].Side:
			return false
		case sliceFMSC[i].Count < sliceFMSC[j].Count:
			return true
		case sliceFMSC[i].Count > sliceFMSC[j].Count:
			return false
		default:
			return false
		}
	})
}

// molecular가 유사한지를 판단한다. 최대 변형 개수 이하여야 한다.
func isLikely(ms1, ms2 string, maxEndure int) bool {
	msLen := len(ms1)
	if msLen != len(ms2) {
		return false
	}
	differ := 0
	for ix := 0; ix < msLen; ix++ {
		if ms1[ix] != ms2[ix] {
			differ++
		}
	}
	return differ <= maxEndure
}

// RemoveSimilarMolecular
//
// 1건만 있는 FMSC를 최대 수량의 FMSC에 더하고 해당 FMSC는 제거한다.
//
func RemoveSimilarMolecular(sliceFMSC *[]FMSC) {
	type tFMSC struct {
		ord    int
		s      FMSC
		remove bool
	}

	var ordered []tFMSC
	for ix, v := range *sliceFMSC {
		ordered = append(ordered, tFMSC{ix, v, false})
	}

	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].s.Count < ordered[j].s.Count
	})

	// 한 개만 존재하는 것들을 상위의 유사한 FMS에 추가하고 삭제한다.
	for ix := 0; ix < len(ordered) && ordered[ix].s.Count == 1; ix++ {
		for jx := len(ordered) - 1; jx > ix && ordered[jx].s.Count > 1; jx-- {
			if isLikely(ordered[ix].s.Molecular, ordered[jx].s.Molecular, 1) {
				if ordered[ix].s.Side == ordered[jx].s.Side {
					// 유사한 molecular-side에 추가한다. 그리고 삭제 마크를 한다.
					ordered[jx].s.Count += ordered[ix].s.Count
					ordered[ix].remove = true
					break
				}
			}
		}
	}

	// 원래 slice를 재건한다.
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].ord < ordered[j].ord
	})
	*sliceFMSC = nil
	for _, v := range ordered {
		if !v.remove {
			*sliceFMSC = append(*sliceFMSC, v.s)
		}
	}
}
