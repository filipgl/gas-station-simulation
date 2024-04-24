package stages

import (
	"sync"
)

type Group struct {
	wg      *sync.WaitGroup
	inbound <-chan *Car
	stages  []*Server
}

func (g *Group) Serve() <-chan *Car {
	var outbound <-chan *Car

	dispatcher := DispatcherFactory(g.wg, g.inbound, g.stages)
	dispatcher.Dispatch()

	outbounds := make([]<-chan *Car, 0, len(g.stages))

	for _, stage := range g.stages {
		outbounds = append(outbounds, (*stage).Serve())
	}

	merger := MergerFactory(g.wg, outbounds)
	outbound = merger.Merge()

	return outbound
}

func GroupFactory(wg *sync.WaitGroup, inbound <-chan *Car, stages []*Server) *Group {
	return &Group{
		wg,
		inbound,
		stages,
	}
}
