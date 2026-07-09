package cfg

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/zenpaw-labs/skypaw/utils"
)

const (
	Metric   UnitSystem = "Metric"
	Imperial UnitSystem = "Imperial"
)

type UnitSystem string

type UserConfig struct {
	UserCity                      string     `json:"city"`
	OptionalLocationProvider      int        `json:"location_provider_id"`
	WindowsLocalLocationDetection bool       `json:"windows_local_location_detection"`
	HideSunBar                    bool       `json:"hide_sun_bar"`
	Units                         UnitSystem `json:"units"`
	// TODO: Hints
	ShowHints bool `json:"show_hints"`
	// TODO: Colorful terminal
	ColorfulTUI bool `json:"colorful_tui"`
}

func LoadConfig() UserConfig {
	cfg := DefaultConfig()
	file, err := utils.GetConfigFile()
	if err != nil {
		cfg = DefaultConfig()
		SaveConfig(cfg)
		return cfg
	}
	data, _ := os.ReadFile(file)
	json.Unmarshal(data, &cfg)
	SaveConfig(cfg)
	slog.Info("Loaded config", "config", cfg)
	return cfg
}

func SaveConfig(cfg UserConfig) error {
	data, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}
	dir := utils.GetConfigDir()
	os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, "config.json")
	return os.WriteFile(path, data, 0644)
}

func (c UserConfig) FormatTemp(celsius float64) string {
	if c.Units == Imperial {
		f := celsius*9/5 + 32
		return fmt.Sprintf("%.1f°F", f)
	}
	return fmt.Sprintf("%.1f°C", celsius)

}

func ParseUnitSystem(s string) UnitSystem {
	switch strings.ToLower(s) {
	case "imperial":
		return Imperial
	default:
		return Metric
	}
}

func DefaultConfig() UserConfig {
	return UserConfig{
		OptionalLocationProvider:      3,
		WindowsLocalLocationDetection: true,
		HideSunBar:                    false,
		Units:                         Metric,
		ShowHints:                     false,
	}
}
