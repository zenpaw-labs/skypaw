package utils

import (
	"encoding/json"
	"os"
)

// TODO: Config file

const (
	Metric UnitSystem = iota
	Imperial
)

type UnitSystem int

type UserConfig struct {
	userCity         string
	optionalProvider int
	hideSunBar       bool
	unitSystem       UnitSystem
}

func LoadConfig() UserConfig {
	cfg := UserConfig{}
	file, err := GetConfigFile()
	if err != nil {
		cfg = DefaultConfig()
		SaveConfig(cfg)
		return cfg
	}
	data, _ := os.ReadFile(file)
	json.Unmarshal(data, &cfg)
	return cfg
}

func SaveConfig(cfg UserConfig) error {
	data, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}
	path, _ := GetConfigFile()
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func DefaultConfig() UserConfig {
	return UserConfig{
		optionalProvider: 2,
		hideSunBar:       false,
		unitSystem:       Metric,
	}
}
