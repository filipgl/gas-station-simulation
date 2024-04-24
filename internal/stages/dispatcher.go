package stages

import (
	"sync"
)

type Dispatcher struct {
	wg      *sync.WaitGroup
	inbound <-chan *Car
	stages  []*Server
}

func DispatcherFactory(wg *sync.WaitGroup, inbound <-chan *Car, stages []*Server) *Dispatcher {
	return &Dispatcher{wg, inbound, stages}
}

func (d *Dispatcher) Dispatch() {
	d.wg.Add(1)
	go func() {
		for car := range d.inbound {
			placed := false
			for !placed {
				best := *FindMin(d.stages)
				placed = best.AddCar(car)
			}
		}
		for _, stage := range d.stages {
			(*stage).CloseInbound()
		}
		d.wg.Done()
	}()
}
