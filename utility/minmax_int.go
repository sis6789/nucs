package utility

var MaxIntValue = int(^uint(0) >> 1)
var MinIntValue = -MaxIntValue - 1

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
