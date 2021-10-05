package merge_values

import "log"

//  merge return values as interface{}

type GoSubResult struct {
	wChan chan interface{}
	V     []interface{}
}

func New() *GoSubResult {
	var newResult GoSubResult
	newResult.wChan = make(chan interface{})
	// start receiver, repeat until closed
	go func() {
		for v := range newResult.wChan {
			newResult.V = append(newResult.V, v)
		}
	}()
	return &newResult
}

// Close - Terminate result merging go-routine
func (x *GoSubResult) Close() {
	defer func() {
		v := recover()
		if v != nil {
			log.Print("close on closed channel")
		}
	}()
	close(x.wChan)
}

// Send - send result
func (x *GoSubResult) Send(v interface{}) {
	defer func() {
		v := recover()
		if v != nil {
			log.Print("send on closed channel")
		}
	}()
	x.wChan <- v
}
