package utility

import (
	"fmt"
	"testing"
)

func TestMarkNonAcgt(t *testing.T) {
	fmt.Println(MarkNonAcgt("ACGTACGTACGTACGTACGT***ACGT121521625"))
	fmt.Println(ShowOnlyDiffer("...AA..xyz*..."))
}
