package limitGoSub

type LimitGoSub struct {
	waitChannel chan struct{}
}

func New(limit int) *LimitGoSub {
	var wBlock LimitGoSub
	wBlock.waitChannel = make(chan struct{}, limit)
	return &wBlock
}

func (x *LimitGoSub) Close() {
	defer func() {
		_ = recover()
	}()
	close(x.waitChannel)
}

func (x *LimitGoSub) Wait() {
	defer func() {
		_ = recover()
	}()
	// wait until channel write success
	x.waitChannel <- struct{}{}
	return
}

func (x *LimitGoSub) Done() {
	defer func() {
		_ = recover()
	}()
	// make room at waitChannel
	<-x.waitChannel
	return
}

func (x *LimitGoSub) Count() int {
	defer func() {
		_ = recover()
	}()
	return len(x.waitChannel)
}
