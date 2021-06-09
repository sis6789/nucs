package utility

import (
	"fmt"
	"testing"
)

func TestMinMax(t *testing.T) {
	fmt.Println(MinInt(-223, 3, 4, 5, 6, 7))
	fmt.Println(MaxInt(3, 4, 5, 6, 7, -875489))
}
