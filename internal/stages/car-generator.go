package stages

import (
	"gas-station/internal/stats"
	"gas-station/internal/utils"
	"sync"
	"time"
)

type CarGenerator struct {
	wg               *sync.WaitGroup
	toGenerate       int
	lineCapacity     int
	durationInterval [2]time.Duration
	types            []string
	stats            *stats.Stats
}

func (cg *CarGenerator) Generate() <-chan *Car {
	outbound := make(chan *Car, cg.lineCapacity)

	cg.wg.Add(1)
	go func() {
		for i := 0; i < cg.toGenerate; i++ {
			car := &Car{Type: cg.types[utils.RandomIntFromRange(0, len(cg.types)-1)]}
			select {
			case outbound <- car:
			default:
				cg.stats.NotServed++
			}

			time.Sleep(utils.RandomDurationFromRange(cg.durationInterval[0], cg.durationInterval[1]))
		}
		close(outbound)
		cg.wg.Done()
	}()

	return outbound
}

func CarGeneratorFactory(
	wg *sync.WaitGroup,
	numberOfCars int,
	lineCapacity int,
	durationInterval [2]time.Duration,
	types []string,
	stats *stats.Stats,
) *CarGenerator {
	return &CarGenerator{
		wg,
		numberOfCars,
		lineCapacity,
		durationInterval,
		types,
		stats,
	}
}
