package read_fastq_both

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test_All(t *testing.T) {
	w := PairList(`D:\keyomics_test\2020MG0213\etc\HN00122843`, `^[\w_-]+\.(fastq|fq)$`)
	var wg sync.WaitGroup
	for _, f2 := range w {
		{
			z := New()
			z.OpenPair(`D:\keyomics_test\2020MG0213\etc\HN00122843`, strings.Split(f2, ";")[0])
		}
		wg.Add(1)
		go func(pP string, pF2 string) {
			defer wg.Done()
			x := New()
			start := time.Now()
			f2t := strings.Split(pF2, ";")
			x.Open(pP, f2t...)
			if !x.AtRec(3000) {
				fmt.Println("EOF before", math.MaxInt64)
			}
			fmt.Println(x)
			if !x.AtRec(1969553) {
				fmt.Println("EOF before", math.MaxInt64)
			}
			fmt.Println(x)
			if !x.AtRec(1969554) {
				fmt.Println("EOF before", math.MaxInt64)
			}
			fmt.Println(x)
			if !x.AtRec(math.MaxInt32) {
				fmt.Println("EOF before", math.MaxInt64)
			}
			fmt.Println(x)
			x.Close()
			since := time.Since(start)
			fmt.Println(since)
		}(`D:\keyomics_test\2020MG0213\etc\HN00122843`, f2)
		wg.Wait()
	}

}
