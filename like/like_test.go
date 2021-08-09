package like

import (
	"fmt"
	"testing"
)

func TestLike(t *testing.T) {
	gotStart, gotRatio, mLen, gotM1Str, gotM2Str := Like(
		"ACGTACGTACGTACGTACGTACGTACGT"+"GATCGGAAGAGCACACGTCTGAACTCCAGTC",
		"GATCGGAAGAGCACACGTCTGAACTCCAGTCAC")
	fmt.Println(gotStart, gotRatio, mLen, gotM1Str, gotM2Str)
}
