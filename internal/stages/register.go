package stages

import (
	"gas-station/internal/stats"
	"gas-station/internal/utils"
	"sync"
	"time"
)

type Register struct {
	wg               *sync.WaitGroup
	inbound          chan *Car
	durationInterval [2]time.Duration
	totalInLine      int
	mutex            sync.Mutex
	stats            *stats.ServerStats
}

func (r *Register) IncrementCounter() {
	r.mutex.Lock()
	r.totalInLine++
	r.mutex.Unlock()
}

func (r *Register) DecrementCounter() {
	r.mutex.Lock()
	r.totalInLine--
	r.mutex.Unlock()
}

func (r *Register) AddCar(car *Car) bool {
	car.RegisterLineEnter = time.Now()

	select {
	case r.inbound <- car:
		r.IncrementCounter()
		return true
	default:
		return false
	}
}

func (r *Register) CloseInbound() {
	close(r.inbound)
}

func (r *Register) GetOccupyingCount() int {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.totalInLine
}

func (r *Register) Serve() <-chan *Car {
	outbound := make(chan *Car)

	r.wg.Add(1)
	go func() {
		for car := range r.inbound {
			UpdateStats(r.stats, time.Since(car.RegisterLineEnter))

			time.Sleep(utils.RandomDurationFromRange(r.durationInterval[0], r.durationInterval[1]))
			car.HasPaid = true
			outbound <- car

			r.DecrementCounter()
		}
		close(outbound)
		r.wg.Done()
	}()

	return outbound
}

func RegisterFactory(
	wg *sync.WaitGroup,
	durationInterval [2]time.Duration,
	lineCapacity int,
	stats *stats.ServerStats,
) *Register {
	return &Register{
		wg:               wg,
		inbound:          make(chan *Car, lineCapacity),
		durationInterval: durationInterval,
		stats:            stats,
	}
}

func RegisterGroupFactory(
	wg *sync.WaitGroup,
	inbound <-chan *Car,
	durationInterval [2]time.Duration,
	lineCapacity int,
	count int,
	stats *stats.ServerStats,
) *Group {
	registers := make([]*Server, 0, count)

	for i := 0; i < count; i++ {
		register := Server(RegisterFactory(wg, durationInterval, lineCapacity, stats))
		registers = append(registers, &register)
	}

	return GroupFactory(wg, inbound, registers)
}
