package main

import (
	"fmt"
	"gas-station/internal/config"
	"gas-station/internal/simulation"
	"gas-station/internal/stats"
	"time"
)

func main() {
	conf, configErr := config.ReadConfig("simulation-config.yaml")

	if configErr != nil {
		fmt.Println("Unable to get configuration")
		panic(configErr)
	}

	start := time.Now()
	statistics, simErr := simulation.RunSimulation(conf)

	if simErr != nil {
		fmt.Println("Error occurred during simulation")
		panic(simErr)
	}

	fmt.Printf("Simulation took %v\n ", time.Since(start))

	saveErr := stats.WriteStats("simulation-stats.yaml", statistics)

	if saveErr != nil {
		fmt.Println("Statistics can not be saved")
		panic(saveErr)
	}
}
