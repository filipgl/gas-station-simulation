package stages

import (
	"gas-station/internal/stats"
	"time"
)

type Car struct {
	Type              string
	HasPaid           bool
	PumpLineEnter     time.Time
	RegisterLineEnter time.Time
}

func UpdateStats(stats *stats.ServerStats, duration time.Duration) {
	stats.Mutex.Lock()
	if stats.TotalCars == 0 {
		stats.AvgQueueTime = duration
	} else {
		avgFloat := float64(stats.AvgQueueTime.Nanoseconds())
		durationFloat := float64(duration.Nanoseconds())
		countFloat := float64(stats.TotalCars)

		stats.AvgQueueTime = time.Duration(
			countFloat/(countFloat+1)*avgFloat + durationFloat/(countFloat+1),
		)
	}

	if duration > stats.MaxQueueTime {
		stats.MaxQueueTime = duration
	}

	stats.TotalTime += duration
	stats.TotalCars++
	stats.Mutex.Unlock()
}
