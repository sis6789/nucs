package utility

import "strings"

func ValueLine(outer1, outer2, inner1, inner2 int, displayStr, fill string) string {
	fillLeft := fill[0:1]
	fillRight := fillLeft
	fillCenter := displayStr
	if len(fill) >= 2 {
		fillRight = fill[1:2]
	}
	// handle too bing inner string
	if inner1 < outer1 {
		displayStr = displayStr[outer1-inner1:]
		fillCenter = displayStr
		inner1 = outer1
	}
	if inner1+len(displayStr)-1 > outer2 {
		displayStr = displayStr[:len(displayStr)-((inner1+len(displayStr)-1)-outer2)]
		fillCenter = displayStr
		inner2 = outer2
	}

	if len(displayStr) < (inner2 - inner1 + 1) {
		fillCenter += strings.Repeat("_", (inner2-inner1+1)-len(displayStr))
	} else if len(displayStr) > (inner2 - inner1 + 1) {
		fillCenter = displayStr[0 : inner2-inner1+1]
	}
	leftAppend := strings.Repeat(fillLeft, inner1-outer1)
	rightAppend := strings.Repeat(fillRight, outer2-(inner1+len(displayStr)-1))
	return leftAppend + fillCenter + rightAppend
}
