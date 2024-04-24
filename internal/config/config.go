package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"time"
)

type CarsConfig struct {
	Count          int
	LineCapacity   int           `yaml:"line_capacity"`
	ArrivalTimeMin time.Duration `yaml:"arrival_time_min"`
	ArrivalTimeMax time.Duration `yaml:"arrival_time_max"`
}

type StationConfig struct {
	Count        int
	LineCapacity int           `yaml:"line_capacity"`
	ServeTimeMin time.Duration `yaml:"serve_time_min"`
	ServeTimeMax time.Duration `yaml:"serve_time_max"`
}

type RegisterConfig struct {
	Count         int
	LineCapacity  int           `yaml:"line_capacity"`
	HandleTimeMin time.Duration `yaml:"handle_time_min"`
	HandleTimeMax time.Duration `yaml:"handle_time_max"`
}

type StationsConfig struct {
	Gas      StationConfig
	LPG      StationConfig `yaml:"lpg"`
	Diesel   StationConfig
	Electric StationConfig
}

type Config struct {
	Cars         CarsConfig
	Stations     StationsConfig
	StationTypes []string
	Registers    RegisterConfig
}

func ReadConfig(path string) (*Config, error) {
	config := &Config{}
	yamlFile, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, config)

	if err != nil {
		return nil, err
	}

	for _, sf := range reflect.VisibleFields(reflect.TypeOf(config.Stations)) {
		config.StationTypes = append(config.StationTypes, sf.Name)
	}

	return config, nil
}

func GetStationConfig(configuration *Config, typeName string) StationConfig {
	value := reflect.ValueOf(configuration.Stations)
	anyType := reflect.Indirect(value).FieldByName(typeName).Interface()
	return anyType.(StationConfig)
}
