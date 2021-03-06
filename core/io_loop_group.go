package core

import "gunplan.top/concurrentNet/buffer"

func NewIOLoopGroup() *ioLoopGroup {
	return &ioLoopGroup{}
}

type ioLoopGroup struct {
	loops []*ioLoop
	index int
}

func (g *ioLoopGroup) registe(lp *ioLoop) {
	g.loops = append(g.loops, lp)
}

func (g *ioLoopGroup) create(sum int) error {
	for i := 0; i < sum; i++ {
		lp, err := NewIOLoop(i, buffer.NewLikedBufferAllocator())
		if err != nil {
			return err
		}
		g.registe(lp)
	}
	return nil
}

//for put new accept connection in ioLoop load balance
func (g *ioLoopGroup) next() *ioLoop {
	g.index++
	size := len(g.loops)
	if g.index >= size {
		g.index -= size
	}
	return g.loops[g.index]
}

func (g *ioLoopGroup) iterate(f func(*ioLoop) bool) {
	for _, loop := range g.loops {
		if !f(loop) {
			break
		}
	}
}
