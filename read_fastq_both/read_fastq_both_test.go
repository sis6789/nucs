package read_fastq_both

import (
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test_All(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//w := PairList(`T:/data/2022MG0804_HN176113_LIVO`, `(?i)(?P<gname>LIVO)(?P<gnum>\d+)-(?P<fnum>\d+)[._-]R?(?P<rnum>\d+)\.(?P<ext>fastq|fq)(?P<gzip>\.gz)?`)
	w := PairList(`M:\tmp\`, `(?i)(?P<gname>LIVO)(?P<gnum>\d+)-(?P<fnum>\d+)[._-]R?(?P<rnum>\d+)\.(?P<ext>fastq|fq)(?P<gzip>\.gz)?`)
	var wg sync.WaitGroup
	limit := 0
	for _, f2 := range w {
		limit++
		if limit >= 2 {
			break
		}
		wg.Add(1)
		go func(pP string, pF2 string) {
			defer wg.Done()
			x := New()
			start := time.Now()
			f2t := strings.Split(pF2, ";")
			x.Open(pP, f2t...)
			if !x.AtRec(30) {
				fmt.Println("EOF before", math.MaxInt64)
			}
			fmt.Println(x.Text[0])
			fmt.Println(x.Text[1])
			//if !x.AtRec(1969553) {
			//	fmt.Println("EOF before", math.MaxInt64)
			//}
			//fmt.Println(x.recCount)
			if !x.AtRec(40) {
				fmt.Println("EOF before", math.MaxInt64)
			}
			fmt.Println(x.Text[0])
			fmt.Println(x.Text[1])
			if !x.AtRec(math.MaxInt32) {
				fmt.Println("EOF before", math.MaxInt64)
			}
			fmt.Println(x.recCount)
			x.Close()
			since := time.Since(start)
			fmt.Println(since)
		}(`M:\tmp\`, f2)
		wg.Wait()
	}

}
