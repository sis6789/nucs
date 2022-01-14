package like

import (
	"github.com/biogo/biogo/align"
	"github.com/biogo/biogo/alphabet"
	"github.com/biogo/biogo/feat"
	"github.com/biogo/biogo/seq/linear"
)

var alphaWithBlank = alphabet.MustComplement(alphabet.NewComplementor(
	"-acmgrsvtwyhkdbn *",
	feat.DNA,
	alphabet.MustPair(alphabet.NewPairing("acmgrsvtwyhkdbnxACMGRSVTWYHKDBNX- *", "tgkcysbawrdmhvnxTGKCYSBAWRDMHVNX- *")),
	'-', 'n',
	false,
))

var smith = align.SW{ // special for blank handling
	{0, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, 2, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, 2, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, 2, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, 2, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, 2, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, 2, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, 2, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, 2, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, 2, -1, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 2, -1, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 2, -1, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 2, -1, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 2, -1, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 2, -1, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 2, -1, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 2, -1, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 2, -1},
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 2},
}

//var splitExpStr = `\[(\d+),(\d+)\)/\[(\d+),(\d+)\)=(-?\d+)|\[(\d+),(\d+)\)/(-)=(-?\d+)|(-)/\[(\d+),(\d+)\)=(-?\d+)`
//var splitExp = regexp.MustCompile(splitExpStr)

type scorer interface{ Score() int }
type Match struct {
	S1    int
	E1    int
	S2    int
	E2    int
	Score int
}

func Like(s1, s2 []byte) (Match, error) {
	// use smith waterman
	fsa := &linear.Seq{Seq: alphabet.BytesToLetters(s1)}

	fsa.Alpha = alphaWithBlank
	fsb := &linear.Seq{Seq: alphabet.BytesToLetters(s2)}
	fsb.Alpha = alphaWithBlank // alphabet.DNAgapped

	alignList, err := smith.Align(fsa, fsb)
	if err != nil {
		return Match{}, err
	}
	highestMatch := Match{Score: -1000}
	for _, aln := range alignList {
		score := aln.(scorer).Score()
		s1 := aln.Features()[0].Start()
		e1 := aln.Features()[0].End()
		s2 := aln.Features()[1].Start()
		e2 := aln.Features()[1].End()
		if highestMatch.Score < score {
			highestMatch = Match{
				S1:    s1,
				E1:    e1,
				S2:    s2,
				E2:    e2,
				Score: score,
			}
		}
	}
	return highestMatch, nil
}
