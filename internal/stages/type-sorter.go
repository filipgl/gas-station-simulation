package stages

import "sync"

type TypeSorter struct {
	wg      *sync.WaitGroup
	inbound <-chan *Car
	types   []string
}

func (ts *TypeSorter) Sort() map[string]chan *Car {
	outboundsMap := make(map[string]chan *Car)

	for _, typeName := range ts.types {
		outboundsMap[typeName] = make(chan *Car)
	}

	ts.wg.Add(1)
	go func() {
		for car := range ts.inbound {
			outboundsMap[car.Type] <- car
		}
		for _, channel := range outboundsMap {
			close(channel)
		}
		ts.wg.Done()
	}()

	return outboundsMap
}

func TypeSorterFactory(wg *sync.WaitGroup, inbound <-chan *Car, types []string) *TypeSorter {
	return &TypeSorter{wg, inbound, types}
}
