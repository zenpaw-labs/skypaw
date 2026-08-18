package cfg

import (
	"encoding/json"
	"fmt"

	"charm.land/log/v2"

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

const (
	MaxHoursBefore = 24
	MaxHoursAfter  = 48
)

type UserConfig struct {
	UserCity                      string     `json:"city"`
	OptionalLocationProvider      int        `json:"location_provider_id"`
	Units                         UnitSystem `json:"units"`
	DiagramHoursBefore            int        `json:"diagram_hours_before"`
	DiagramHoursAfter             int        `json:"diagram_hours_after"`
	WindowsLocalLocationDetection bool       `json:"windows_local_location_detection"`
	UseWeatherCache               bool       `json:"use_weather_cache"`
	HideDiagram                   bool       `json:"hide_diagram"`
	ShowHints                     bool       `json:"show_hints"`
	ColorfulTUI                   bool       `json:"colorful_tui"`
	AlwaysRunDebugger             bool       `json:"always_run_debugger"`
}

func LoadConfig() UserConfig {
	cfg := DefaultConfig()
	file, err := utils.GetConfigFile()
	if err != nil {
		cfg = DefaultConfig()
		SaveConfig(cfg)
		log.Debug("Default config used and created successfully.", "config", cfg)
		return cfg
	}
	data, _ := os.ReadFile(file)
	json.Unmarshal(data, &cfg)
	SaveConfig(cfg)
	cfg.Validate()
	log.Debug("Config loaded successfully.", "config", cfg)
	return cfg
}

func SaveConfig(cfg UserConfig) error {
	data, err := json.MarshalIndent(cfg, "", "    ")
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
		HideDiagram:                   false,
		Units:                         Metric,
		ShowHints:                     true,
		ColorfulTUI:                   true,
		DiagramHoursBefore:            2,
		DiagramHoursAfter:             6,
		UseWeatherCache:               true,
	}
}

func (c *UserConfig) Validate() {
	example := DefaultConfig()

	c.Units = ParseUnitSystem(string(c.Units))

	if c.DiagramHoursAfter < 0 || c.DiagramHoursAfter > MaxHoursAfter {
		c.DiagramHoursAfter = example.DiagramHoursAfter
	}

	if c.DiagramHoursBefore < 0 || c.DiagramHoursBefore > MaxHoursBefore {
		c.DiagramHoursBefore = example.DiagramHoursBefore
	}
}
