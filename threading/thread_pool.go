package threading

import (
	"sync"
	"time"
)

type StandThreadPool interface {
	ThreadPool
	Init(core, ext int, d time.Duration, w int, strategy func(interface{})) StandThreadPool
}

type threadPoolImpl struct {
	BasePool
	s              func(interface{})
	core           int
	workQueue      chan *Task
	controlChannel chan ControlType
	g              sync.WaitGroup
	w              sync.Mutex
}

func (t *threadPoolImpl) Init(core, ext int, span time.Duration, w int, strategy func(interface{})) StandThreadPool {
	t.BasePool.Init(t.addQueue0)
	t.core = core
	t.workQueue = make(chan *Task, w)
	t.controlChannel = make(chan ControlType, t.core)
	t.s = strategy
	t.g.Add(t.core)
	return t
}

func (t *threadPoolImpl) Boot() Launcher {
	t.status.whenThreadBooting()
	for i := 0; i < t.core; i++ {
		go t.LaunchWork()
	}
	t.status.whenThreadBooted()
	return t
}

func (t *threadPoolImpl) LaunchWork() {
	for {
		select {
		case task, ok := <-t.workQueue:
			// core execute unit
			if ok {
				task.rev <- task.t(task.param...)
			}
			// control unit
		case op := <-t.controlChannel:
			switch op {
			case SHUTDOWN:
				if len(t.workQueue) != 0 {
					t.consumeRemain()
				}
				fallthrough

			case ShutdownNow:
				t.controlChannel <- op
				t.g.Done()
				return
			}
		}
	}
}

func (t *threadPoolImpl) consumeRemain() {
	for len(t.workQueue) != 0 {
		// The necessity og locking here is that
		// we have to make sure operator get length
		// and operator consume the channel is an
		// atomic operation.
		t.w.Lock()
		if len(t.workQueue) != 0 {
			task := <-t.workQueue
			t.w.Unlock()
			task.rev <- task.t(task.param...)
		} else {
			t.w.Unlock()
		}
	}
}

func (t *threadPoolImpl) LaunchWorkExt() {
	t.g.Add(1)
	t.LaunchWork()
}

func (t *threadPoolImpl) Join() {
	t.g.Wait()
}

func (t *threadPoolImpl) Shutdown() {
	t.waitStop(SHUTDOWN)
}

func (t *threadPoolImpl) ShutdownNow() {
	close(t.workQueue)
	t.waitStop(ShutdownNow)
}

func (t *threadPoolImpl) waitStop(c ControlType) {
	t.status.whenThreadStopping()
	t.controlChannel <- c
	go func() {
		t.g.Wait()
		close(t.controlChannel)
		t.status.whenThreadStopped()
	}()
}

func (t *threadPoolImpl) addQueue0(task *Task) {
	t.workQueue <- task
}
