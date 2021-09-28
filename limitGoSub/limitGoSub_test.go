package limitGoSub

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLimitGoSub_Close(t *testing.T) {
	x := New(2)
	start := time.Now()
	go func() {
		for {
			time.Sleep(1 * time.Second)
			x.Done()
		}
	}()
	var wg sync.WaitGroup
	for j := 0; j < 20; j++ {
		x.Wait()
		wg.Add(1)
		go func(ix int) {
			fmt.Println("ZZZ", ix, time.Since(start))
			wg.Done()
		}(j)
	}
	wg.Wait()
	fmt.Println("EOJ")
}
