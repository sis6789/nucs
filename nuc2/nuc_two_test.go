package nuc2

import (
	"fmt"
	"testing"
)

func TestNuc2(t *testing.T) {
	fmt.Println(string(Nuc2D('A', 'C')))
	fmt.Println(Nuc2DString("A   ", "ACGT"))
	fmt.Println(Nuc2DString("CCCC", "ACGT"))
	fmt.Println(Nuc2DString("GGGG", "ACGT"))
	fmt.Println(Nuc2DString("TTTT", "ACGT"))
}
