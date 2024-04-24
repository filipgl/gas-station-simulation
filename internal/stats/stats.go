package stats

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
	"time"
)

type ServerStats struct {
	TotalCars    int           `yaml:"total_cars"`
	TotalTime    time.Duration `yaml:"total_time"`
	AvgQueueTime time.Duration `yaml:"avg_queue_time"`
	MaxQueueTime time.Duration `yaml:"max_queue_time"`
	Mutex        sync.Mutex    `yaml:"-"`
}

type Stats struct {
	Gas       ServerStats
	LPG       ServerStats `yaml:"lpg"`
	Diesel    ServerStats
	Electric  ServerStats
	Registers ServerStats
	NotServed int `yaml:"not_served"`
}

func WriteStats(path string, stats *Stats) error {
	yamlContent, err := yaml.Marshal(stats)

	if err != nil {
		return err
	}

	err = os.WriteFile(path, yamlContent, 0630)

	if err != nil {
		return err
	}

	return nil
}

func GetServerStats(stats *Stats, typeName string) (*ServerStats, error) {
	// Didn't find solution using reflect

	switch typeName {
	case "LPG":
		return &stats.LPG, nil
	case "Diesel":
		return &stats.Diesel, nil
	case "Electric":
		return &stats.Electric, nil
	case "Gas":
		return &stats.Gas, nil
	}

	return nil, errors.New("unexpected server type")
}
