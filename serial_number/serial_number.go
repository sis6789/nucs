package serial_number

type SerialNumber struct {
	int32Channel chan int32
}

func New() *SerialNumber {
	var sn SerialNumber
	sn.int32Channel = make(chan int32, 5)
	go sn.generator()
	return &sn
}

func (x *SerialNumber) generator() {
	defer func() {
		recover() // disable chan panic
	}()
	var num = int32(0)
	for {
		num++
		x.int32Channel <- num
	}
}

func (x *SerialNumber) Next() int32 {
	return <-x.int32Channel
}

func (x *SerialNumber) Close() {
	close(x.int32Channel)
}
