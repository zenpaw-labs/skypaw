package cfg

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"

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

var hexRe = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)

type UnitSystem string

const (
	MaxHoursBefore = 24
	MaxHoursAfter  = 48
)

type UserConfig struct {
	UserCity                      string       `json:"city"`
	OptionalLocationProvider      int          `json:"location_provider_id"`
	Units                         UnitSystem   `json:"units"`
	DiagramHoursBefore            int          `json:"diagram_hours_before"`
	DiagramHoursAfter             int          `json:"diagram_hours_after"`
	WindowsLocalLocationDetection bool         `json:"windows_local_location_detection"`
	UseWeatherCache               bool         `json:"use_weather_cache"`
	SingleLineOutput              bool         `json:"single_line_output"`
	HideDiagram                   bool         `json:"hide_diagram"`
	ShowHints                     bool         `json:"show_hints"`
	ColorfulTUI                   bool         `json:"colorful_tui"`
	BadWeatherNotification        bool         `json:"bad_weather_notification"`
	BadWeatherLookahead           int          `json:"bad_weather_lookahead"`
	Template                      []string     `json:"template"`
	SingleLineTemplate            string       `json:"single_line_template"`
	Colors                        ColorsConfig `json:"colors"`
	AlwaysRunDebugger             bool         `json:"always_run_debugger"`
}

var DefaultTemplate = []string{
	`{{printf "📍 %s" .Location | theme "location"}}`,
	``,
	`{{.Art}}`,
	``,
	`{{.Temp | theme "temp"}}`,
	`{{.Condition | theme "weather"}}`,
	`{{.WindDirection | theme "weather"}} {{.WindSpeed | theme "weather"}}`,
	``,
	`{{.TimeSec | theme "time"}}`,
	`{{.DateLong | theme "date"}}{{if .Offline}}`,
	`{{.OfflineNotice}}{{end}}{{if .Hints}}`,
	``,
	`{{.Hints}}{{end}}`,
}

const DefaultSingleLineTemplate = `{{printf "📍 %s" .Location | theme "location"}}: {{.Temp | theme "temp"}}, {{.Condition | theme "weather"}}{{if .Offline}} {{"(offline)" | theme "offline"}}{{end}}`

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
	cfg.Validate()
	SaveConfig(cfg)
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

func (c UserConfig) ActiveTemplate() string {
	if c.SingleLineOutput {
		return c.SingleLineTemplate
	}
	return strings.Join(c.Template, "\n")
}

func (c UserConfig) DefaultActiveTemplate() string {
	d := DefaultConfig()
	d.SingleLineOutput = c.SingleLineOutput
	return d.ActiveTemplate()
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
		SingleLineOutput:              false,
		BadWeatherNotification:        true,
		BadWeatherLookahead:           6,
		Colors:                        DefaultColors(),
		Template:                      DefaultTemplate,
		SingleLineTemplate:            DefaultSingleLineTemplate,
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

	v := reflect.ValueOf(&c.Colors).Elem()
	d := reflect.ValueOf(DefaultColors())
	for i := 0; i < v.NumField(); i++ {
		if !hexRe.MatchString(v.Field(i).String()) {
			log.Warn("Invalid color, using default", "key", v.Type().Field(i).Name)
			v.Field(i).SetString(d.Field(i).String())
		}
	}

	if len(c.Template) == 0 {
		c.Template = DefaultTemplate
	}
	if strings.TrimSpace(c.SingleLineTemplate) == "" {
		c.SingleLineTemplate = DefaultSingleLineTemplate
	}
}
