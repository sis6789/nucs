package serial_number

import (
	"fmt"
	"log"
)

type SerialNumber struct {
	int32Channel chan int32
	active       bool
}

func New() *SerialNumber {
	var sn SerialNumber
	sn.int32Channel = make(chan int32, 5)
	sn.active = true
	go sn.generator()
	return &sn
}

func (x *SerialNumber) generator() {
	defer func() {
		recover()        // disable chan panic
		x.active = false // set NA
		fmt.Println("EOG")
	}()
	var num = int32(0)
	for {
		num++
		x.int32Channel <- num
	}
}

func (x *SerialNumber) Next() int32 {
	if x.active {
		return <-x.int32Channel
	} else {
		log.Panicln("call after close.")
		return -1
	}
}

func (x *SerialNumber) Close() {
	close(x.int32Channel)
}
