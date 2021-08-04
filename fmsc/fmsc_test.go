package fmsc

import (
	"fmt"
	"testing"
)

func TestSortFMS(t *testing.T) {
	v := []FMSC{
		NewFMSC("1", "ACGTACGT", "1", 10),
		NewFMSC("1", "ACGTACG1", "1", 1),
		NewFMSC("1", "ACGTACG2", "1", 1),
		NewFMSC("1", "ACGTACG3", "1", 1),
		NewFMSC("1", "ACGTACG4", "1", 1),
		NewFMSC("1", "BCGTACGT", "1", 10),
		NewFMSC("1", "BCGTACG1", "1", 1),
		NewFMSC("1", "BCGTACG2", "1", 1),
		NewFMSC("1", "BCGTACG3", "1", 1),
		NewFMSC("1", "BCGTACG4", "1", 1),
		NewFMSC("1", "BCGTACGT", "1", 10),
		NewFMSC("1", "BCGTACX1", "1", 1),
		NewFMSC("1", "BCGTACX2", "1", 1),
		NewFMSC("9", "BCGTACX3", "1", 1),
		NewFMSC("1", "BCGTACX4", "1", 1),
	}

	fmt.Println(v)
	RemoveSimilarMolecular(&v)
	fmt.Println(v)

}
