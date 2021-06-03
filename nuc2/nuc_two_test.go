package nuc2

import (
	"fmt"
	"testing"
)

func TestNuc2(t *testing.T) {
	fmt.Println(string(Nuc2D('A', 'C')))
	fmt.Println(Nuc2DString("A", "ACGT"))
	fmt.Println(Nuc2DString("CCCCAAAAAAAAAAAAAA", "ACGTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"))
	fmt.Println(Nuc2DString("GGGG", "ACGT"))
	fmt.Println(Nuc2DString("TTTT", "ACGT"))
	fmt.Println(Nuc2String("A", "ACGT"))
	fmt.Println(Nuc2String("CCCCAAAAAAAAAAA", "ACGTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"))
	fmt.Println(Nuc2String("GGGG", "ACGT"))
	fmt.Println(Nuc2String("TTTT&*&*&*", "ACGT"))
}
