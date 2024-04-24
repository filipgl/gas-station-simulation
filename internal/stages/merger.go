package stages

import (
	"sync"
)

type Merger struct {
	wg *sync.WaitGroup
	cs []<-chan *Car
}

func MergerFactory(wg *sync.WaitGroup, cs []<-chan *Car) *Merger {
	return &Merger{wg, cs}
}

func (m *Merger) Merge() <-chan *Car {
	var wg2 sync.WaitGroup
	outbound := make(chan *Car)

	output := func(c <-chan *Car) {
		for n := range c {
			outbound <- n
		}
		wg2.Done()
	}

	wg2.Add(len(m.cs))
	for _, c := range m.cs {
		go output(c)
	}

	m.wg.Add(1)
	go func() {
		wg2.Wait()
		close(outbound)
		m.wg.Done()
	}()

	return outbound
}
