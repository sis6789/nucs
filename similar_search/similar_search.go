package similar_search

import (
	"fmt"

	"github.com/biogo/biogo/align"
	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/seq/linear"
)

func Similar(v1, v2 string) (from int, ratio float64, pattern string) {
	fsa := &linear.Seq{Seq: alphabet.BytesToLetters([]byte(v1))}
	fsa.Alpha = alphabet.DNAgapped
	fsb := &linear.Seq{Seq: alphabet.BytesToLetters([]byte(v2))}
	fsb.Alpha = alphabet.DNAgapped

	//		   Query letter
	//  	 -	 A	 C	 G	 T
	// -	 0	-1	-1	-1	-1
	// A	-1	 1	-1	-1	-1
	// C	-1	-1	 1	-1	-1
	// G	-1	-1	-1	 1	-1
	// T	-1	-1	-1	-1	 1
	//
	// Gap open: -5
	fitted := align.FittedAffine{
		Matrix: align.Linear{
			{0, -1, -1, -1, -1},
			{-1, 1, -1, -1, -1},
			{-1, -1, 1, -1, -1},
			{-1, -1, -1, 1, -1},
			{-1, -1, -1, -1, 1},
		},
		GapOpen: -5,
	}

	aln, err := fitted.Align(fsa, fsb)
	if err == nil {
		f2 := aln[0].Features()
		from = f2[0].Start()
		fa := align.Format(fsa, fsb, aln, '-')
		pattern = fmt.Sprint(fa[0])
		fa0 := []byte(pattern)
		fa1 := []byte(fmt.Sprint(fa[1]))
		differ := 0
		for ix := 0; ix < len(fa0); ix++ {
			if fa0[ix] != fa1[ix] {
				differ++
			}
		}
		ratio = 1.0 - float64(differ)/float64(len(fa0))
	} else {
		from = 0
		ratio = 0.0
		pattern = ""
	}
	return from, ratio, pattern
}
