package utility

import (
	"regexp"
	"strings"
)

func MarkNonAcgt(acgtStr string) string {
	patternAcgt := regexp.MustCompile(`([ acgtACGT]*)([^ acgtACGT]*)`)
	tokens := patternAcgt.FindAllStringSubmatch(acgtStr, -1)
	result := ""
	for _, token1 := range tokens {
		result += `<span style='color:black;'>` + token1[1] + `</span>`
		if token1[2] != "" {
			result += `<span style='color:black; background-color: coral;'>` + token1[2] + `</span>`
		}
	}
	return result
}
func ShowOnlyDiffer(acgtStr string) string {
	patternAcgt := regexp.MustCompile(`([ .]*)([^ .]*)`)
	tokens := patternAcgt.FindAllStringSubmatch(acgtStr, -1)
	result := ""
	for _, token1 := range tokens {
		result += `<span style='color:black;'>` + strings.Repeat(" ", len(token1[1])) + `</span>`
		if token1[2] != "" {
			result += `<span style='color:red;'>` + token1[2] + `</span>`
		}
	}
	return result
}
