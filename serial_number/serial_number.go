package serial_number

import (
	"log"
	"sync"
)

type SerialNumber struct {
	intChannel chan int
	active     bool
	mu         sync.Mutex
}

func New() *SerialNumber {
	var sn SerialNumber
	sn.intChannel = make(chan int, 5)
	sn.active = true
	go sn.generator()
	return &sn
}

func (x *SerialNumber) generator() {
	defer func() {
		recover()        // disable chan panic
		x.active = false // set NA
	}()
	var num = int(0)
	for {
		num++
		x.intChannel <- num
	}
}

func (x *SerialNumber) Next() int {
	if x.active {
		x.mu.Lock()
		newSerial := <-x.intChannel
		x.mu.Unlock()
		return newSerial
	} else {
		log.Panicln("call after close.")
		return -1
	}
}

func (x *SerialNumber) Close() {
	close(x.intChannel)
}
