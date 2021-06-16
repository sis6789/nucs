package tail_match

import (
	"fmt"
	"testing"
)

func TestTailMatch(t *testing.T) {
	target := "A123AACGTX12TTTA58439867AXCG4895745984375ACGTX12CGACGTX1"
	{
		isOK, mStart, mLen := MatchAny(target, "ACGTX123", 2, 5)
		fmt.Println(isOK, mStart, mLen)
		if isOK {
			fmt.Println(target[mStart : mStart+mLen])
		}
	}
	{
		isOK, mStart, mLen := MatchTail(target, "ACGTX123", 2, 5)
		fmt.Println(isOK, mStart, mLen)
		if isOK {
			fmt.Println(target[mStart : mStart+mLen])
		}
	}
}
