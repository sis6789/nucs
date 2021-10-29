package serial_number

import (
	"log"
)

type SerialNumber struct {
	intChannel chan int
	active     bool
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
	var num = 0
	for {
		num++
		x.intChannel <- num
	}
}

func (x *SerialNumber) Next() int {
	if x.active {
		return <-x.intChannel
	} else {
		log.Panicln("call after close.")
		return -1
	}
}

func (x *SerialNumber) Close() {
	close(x.intChannel)
}
