package stages

import (
	"gas-station/internal/stats"
	"gas-station/internal/utils"
	"sync"
	"time"
)

type Pump struct {
	wg               *sync.WaitGroup
	inbound          chan *Car
	durationInterval [2]time.Duration
	totalInLine      int
	mutex            sync.Mutex
	stats            *stats.ServerStats
}

func (p *Pump) IncrementCounter() {
	p.mutex.Lock()
	p.totalInLine++
	p.mutex.Unlock()
}

func (p *Pump) DecrementCounter() {
	p.mutex.Lock()
	p.totalInLine--
	p.mutex.Unlock()
}

func (p *Pump) AddCar(car *Car) bool {
	// Placed here to avoid problems when car is taken from inbound too fast
	car.PumpLineEnter = time.Now()

	select {
	case p.inbound <- car:
		p.IncrementCounter()
		// Original place: car.PumpLineEnter = time.Now()
		return true
	default:
		return false
	}
}

func (p *Pump) CloseInbound() {
	close(p.inbound)
}

func (p *Pump) GetOccupyingCount() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.totalInLine
}

func (p *Pump) Serve() <-chan *Car {
	outbound := make(chan *Car)

	p.wg.Add(1)
	go func() {
		for car := range p.inbound {
			UpdateStats(p.stats, time.Since(car.PumpLineEnter))

			time.Sleep(utils.RandomDurationFromRange(p.durationInterval[0], p.durationInterval[1]))
			outbound <- car

			for !car.HasPaid {
			}

			p.DecrementCounter()
		}
		close(outbound)
		p.wg.Done()
	}()

	return outbound
}

func PumpFactory(wg *sync.WaitGroup,
	durationInterval [2]time.Duration,
	lineCapacity int,
	stats *stats.ServerStats,
) *Pump {
	return &Pump{
		wg:               wg,
		inbound:          make(chan *Car, lineCapacity),
		durationInterval: durationInterval,
		stats:            stats,
	}
}

func PumpGroupFactory(wg *sync.WaitGroup,
	inbound <-chan *Car,
	durationInterval [2]time.Duration,
	lineCapacity int,
	count int,
	stats *stats.ServerStats,
) *Group {
	pumps := make([]*Server, 0, count)

	for i := 0; i < count; i++ {
		pump := Server(PumpFactory(wg, durationInterval, lineCapacity, stats))
		pumps = append(pumps, &pump)
	}

	return GroupFactory(wg, inbound, pumps)
}
