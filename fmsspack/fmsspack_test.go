package fmsspack

import (
	"fmt"
	"math"
	"testing"
)

func TestAny(t *testing.T) {
	maxFNum := int(math.Pow(2, 14) - 1)
	{
		fmt.Println(`----------`)
		fmt.Println("max file number is", maxFNum)
		v := Encode(maxFNum, "ACGTG", "", "C")
		fmt.Println(v)
		fmt.Println(v >> 18)
		fileNumber, molecularIndex, shift, strand := Decode(v)
		fmt.Println(fileNumber, molecularIndex, shift, strand)
	}
	{
		fmt.Println(`----------`)
		v := Encode(71, "ACGTG", "", "C")
		fmt.Println(v)
		fmt.Println(v >> 18)
		fileNumber, molecularIndex, shift, strand := Decode(v)
		fmt.Println(fileNumber, molecularIndex, shift, strand)
	}
}
