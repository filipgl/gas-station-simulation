package simulation

import (
	"gas-station/internal/config"
	"gas-station/internal/stages"
	"gas-station/internal/stats"
	"sync"
	"time"
)

func RunSimulation(configuration *config.Config) (*stats.Stats, error) {
	var wg sync.WaitGroup
	statistics := &stats.Stats{}

	gen := stages.CarGeneratorFactory(
		&wg,
		configuration.Cars.Count,
		configuration.Cars.LineCapacity,
		[2]time.Duration{configuration.Cars.ArrivalTimeMin, configuration.Cars.ArrivalTimeMax},
		configuration.StationTypes,
		statistics,
	)

	arrivalLine := gen.Generate()

	sorter := stages.TypeSorterFactory(&wg, arrivalLine, configuration.StationTypes)
	sortedLines := sorter.Sort()

	registerLines := make([]<-chan *stages.Car, 0, len(configuration.StationTypes))

	for stationType, registerLine := range sortedLines {
		stationConf := config.GetStationConfig(configuration, stationType)
		serverStats, err := stats.GetServerStats(statistics, stationType)

		if err != nil {
			return nil, err
		}

		group := stages.PumpGroupFactory(
			&wg,
			registerLine,
			[2]time.Duration{stationConf.ServeTimeMin, stationConf.ServeTimeMax},
			stationConf.LineCapacity,
			stationConf.Count,
			serverStats,
		)

		registerLines = append(registerLines, group.Serve())
	}

	mrg := stages.MergerFactory(&wg, registerLines)
	registerLine := mrg.Merge()

	group := stages.RegisterGroupFactory(
		&wg,
		registerLine,
		[2]time.Duration{configuration.Registers.HandleTimeMin, configuration.Registers.HandleTimeMax},
		configuration.Registers.LineCapacity,
		configuration.Registers.Count,
		&statistics.Registers,
	)

	departureLine := group.Serve()

	// Sink
	wg.Add(1)
	go func() {
		for range departureLine {
		}
		wg.Done()
	}()

	wg.Wait()

	return statistics, nil
}
