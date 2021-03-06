package core

import (
	"log"
	"sync"

	"gunplan.top/concurrentNet/buffer"
	"gunplan.top/concurrentNet/core/netpoll"
)

type ioLoop struct {
	poller   *netpoll.Poller
	channels map[int]Channel
	index    int
	lk       sync.Mutex
	alloc    buffer.Allocator
	l        Pipeline
}

func NewIOLoop(index int, alloc buffer.Allocator) (*ioLoop, error) {
	poller, err := netpoll.NewPoller()
	if err != nil {
		return nil, err
	}
	lp := &ioLoop{
		index:    index,
		poller:   poller,
		channels: make(map[int]Channel),
		alloc:    alloc,
	}
	return lp, nil
}

func (lp *ioLoop) start() {

	if err := lp.poller.Polling(lp.eventHandler); err != nil {
		log.Println(err)
	}
}

func (lp *ioLoop) stop() {
	if err := lp.poller.Trigger(func() error {
		return errLoopShutdown
	}); err != nil {
		log.Printf("index:%d , %v", lp.index, err)
	}
}

func (lp *ioLoop) close() {
	lp.poller.Close()
}

func (lp *ioLoop) eventHandler(fd int, events uint32) error {
	//if channel,ok:=lp.channels[fd];ok{
	//	//switch {
	//	//
	//	//}
	//}
	return nil
}

func (lp *ioLoop) Read(buffer.ByteBuffer, error) {

}
