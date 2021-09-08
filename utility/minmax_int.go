package utility

import "math"

var MaxIntValue = math.MaxInt
var MinIntValue = math.MinInt

func MinInt(x ...int) int {
	lowInt := x[0]
	for _, v := range x {
		if v < lowInt {
			lowInt = v
		}
	}
	return lowInt
}

func MaxInt(x ...int) int {
	highInt := x[0]
	for _, v := range x {
		if v > highInt {
			highInt = v
		}
	}
	return highInt
}
