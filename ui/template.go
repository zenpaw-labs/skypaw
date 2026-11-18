package ui

import (
	"fmt"
	"image/color"
	"strings"
	"text/template"
	"time"

	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"github.com/zenpaw-labs/skypaw/network/weather"
	"github.com/zenpaw-labs/skypaw/utils/cfg"
)

type OutputData struct {
	City          string
	Region        string
	Location      string
	Temp          string
	TempValue     float64
	Condition     string
	Sunrise       string
	Sunset        string
	Time          string
	TimeSec       string
	Date          string
	DateLong      string
	Weekday       string
	Offline       bool
	OfflineNotice string

	Art   string
	Hints string
}

func themeColor(name string) (color.Color, bool) {
	switch name {
	case "location":
		return ColorLocation, true
	case "temp":
		return ColorTemp, true
	case "weather":
		return ColorWeatherName, true
	case "time":
		return ColorTime, true
	case "date":
		return ColorDate, true
	case "hints":
		return ColorHints, true
	case "error":
		return ColorError, true
	case "loading":
		return ColorLoading, true
	case "offline":
		return ColorOfflineOrTimeout, true
	}
	return nil, false
}

func buildTemplate(src string, colorful bool) (*template.Template, error) {
	funcs := template.FuncMap{
		"color": func(c, s string) string {
			if !colorful {
				return s
			}
			return lipgloss.NewStyle().Foreground(lipgloss.Color(c)).Render(s)
		},
		"theme": func(name, s string) string {
			c, ok := themeColor(name)
			if !colorful || !ok {
				return s
			}
			return lipgloss.NewStyle().Foreground(c).Render(s)
		},
		"bold": func(s string) string {
			if !colorful {
				return s
			}
			return lipgloss.NewStyle().Bold(true).Render(s)
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"round": func(f float64) string { return fmt.Sprintf("%.0f", f) },
	}
	return template.New("output").Funcs(funcs).Parse(src)
}

func loadTemplate(c cfg.UserConfig) *template.Template {
	t, err := buildTemplate(c.ActiveTemplate(), c.ColorfulTUI)
	if err == nil {
		return t
	}
	log.Error("Invalid template in config, using default", "err", err)

	t, err = buildTemplate(c.DefaultActiveTemplate(), c.ColorfulTUI)
	if err != nil {
		log.Error("Default template is broken", "err", err)
		return nil
	}
	return t
}

func (m Model) outputData() OutputData {
	loc := m.Location.City
	switch {
	case m.Location.Region != "" && m.Location.City != "":
		loc = m.Location.Region + ", " + m.Location.City
	case m.Location.Region != "":
		loc = m.Location.Region
	}

	unit := "°C"
	if m.Config.Units == cfg.Imperial {
		unit = "°F"
	}
	t := m.Weather.CurrentWeather.Temperature2m

	var offlineNotice string
	if m.isOfflineMode {
		offlineNotice = fmt.Sprintf("Offline mode • Last update at %s",
			m.offlineModeData.LastUpdatedAt.Format("15:04 (02.01.2006)"))
		if m.Config.ColorfulTUI {
			offlineNotice = lipgloss.NewStyle().Foreground(ColorOfflineOrTimeout).Render(offlineNotice)
		}
	}

	return OutputData{
		City:      m.Location.City,
		Region:    m.Location.Region,
		Location:  loc,
		Temp:      fmt.Sprintf("%.1f%s", t, unit),
		TempValue: t,
		Condition: weather.GetCurrentWeatherName(m.Weather.CurrentWeather.WeatherCode),
		Sunrise:   firstTime(m.SunriseAndSunset.Daily.Sunrise),
		Sunset:    firstTime(m.SunriseAndSunset.Daily.Sunset),
		Time:      m.CurrentTime.Format("15:04"),
		TimeSec:   m.CurrentTime.Format("15:04:05"),
		Date:      m.CurrentTime.Format("02.01.2006"),
		DateLong: fmt.Sprintf("%s, %s %d",
			m.CurrentTime.Weekday(), m.CurrentTime.Month(), m.CurrentTime.Day()),
		Weekday:       m.CurrentTime.Weekday().String(),
		Offline:       m.isOfflineMode,
		OfflineNotice: offlineNotice,
	}
}

func firstTime(list []string) string {
	if len(list) == 0 {
		return ""
	}
	t, err := time.Parse("2006-01-02T15:04", list[0])
	if err != nil {
		return ""
	}
	return t.Format("15:04")
}

func (m Model) execTemplate(data OutputData) string {
	if m.tmpl == nil {
		return "template error: see log"
	}
	var sb strings.Builder
	if err := m.tmpl.Execute(&sb, data); err != nil {
		log.Error("Template execute failed", "err", err)
		return "template error: " + err.Error()
	}
	return sb.String()
}

func cropArt(raw string) (string, int) {
	lines := strings.Split(strings.Trim(raw, "\n\r"), "\n")

	minSpaces := -1
	for _, line := range lines {
		line = strings.TrimRight(line, "\r ")
		if len(line) == 0 {
			continue
		}
		leading := len(line) - len(strings.TrimLeft(line, " "))
		if minSpaces == -1 || leading < minSpaces {
			minSpaces = leading
		}
	}
	if minSpaces == -1 {
		minSpaces = 0
	}

	var cropped []string
	artWidth := 0
	for _, line := range lines {
		line = strings.TrimRight(line, "\r ")
		if len(line) >= minSpaces {
			c := line[minSpaces:]
			cropped = append(cropped, c)
			if w := lipgloss.Width(c); w > artWidth {
				artWidth = w
			}
		} else {
			cropped = append(cropped, "")
		}
	}
	for i, line := range cropped {
		if w := lipgloss.Width(line); w < artWidth {
			cropped[i] = line + strings.Repeat(" ", artWidth-w)
		}
	}
	return strings.Join(cropped, "\n"), artWidth
}
