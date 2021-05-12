package applyBtop

import (
	"fmt"
	"testing"
)

func TestAny(t *testing.T) {
	a1 := []string{"-A-C-G-T"}
	fmt.Println(a1)
	Complement(a1)
	fmt.Println(a1)

	bp := "TGGGGAATTTGTTGCAATAAGATCATTGCCTGGAGGTTCCAAGATGGCCGAATAGGAACAGCACACGTCTGAACTCCAGTCACACACTCAGGTCGTTTGGCGTCTTGGGGTTGAAAAAAAAAAA"
	bt := "31-T3C-4"
	req := BtopApplyRequest{bp, 32, 65370217, -40, bt}

	resp := ApplyBtop(req)

	line1 := resp.QueryBP
	spos := resp.GenomeAddress
	qpos := resp.QueryAddress
	modify := resp.ModifyAddress
	d1 := resp.Line1
	d2 := resp.Line2
	dStart := resp.RStart
	dLen := resp.RLen
	fmt.Println(dStart, dLen)
	fmt.Println(d1)
	fmt.Println(d2)
	for ix := 0; ix < len(line1); ix++ {
		fmt.Println(ix+1, string(line1[ix]), spos[ix], qpos[ix], modify[ix])
	}
}
