package ui

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zenpaw-labs/skypaw/ascii"
	"github.com/zenpaw-labs/skypaw/network/geocoding"
	"github.com/zenpaw-labs/skypaw/network/weather"
	"github.com/zenpaw-labs/skypaw/utils/cfg"
)

//TODO: Interactive location picker

type Model struct {
	// Weather
	City             string
	Weather          weather.WeatherResponse
	SunriseAndSunset weather.SunriseAndSunsetResponse
	Location         geocoding.LocationInfo

	// Status
	CurrentTime    time.Time
	CurrentWeekday time.Weekday
	CurrentMonth   time.Month
	IsLoading      int
	Err            error

	// Window
	Width  int
	Height int

	// User config
	Config cfg.UserConfig

	// Other
	Version string
}

func InitialModel(cfg cfg.UserConfig, version string) Model {
	return Model{
		Config:      cfg,
		Version:     version,
		CurrentTime: time.Now(),
		IsLoading:   1,
	}
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	if m.Config.UserCity != "" {
		cmds = append(cmds, FetchLocationByName(m.Config.UserCity))
	} else {
		cmds = append(cmds, FetchLocation(m.Config, m.Config.WindowsLocalLocationDetection))
	}
	cmds = append(cmds, DoTick())
	cmds = append(cmds, DoWeatherRefreshTick())
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case GeocodingMsg:
		m.Location = msg.Data
		if m.IsLoading != 0 {
			m.IsLoading = 2
		}
		return m, FetchWeather(m.Location, m.Config.Units)

	case WeatherMsg:
		m.Weather = msg.Data
		m.Location = msg.LocationInfo
		if m.IsLoading != 0 {
			m.IsLoading = 3
		}
		return m, FetchSunriseAndSunset(m.Location)

	case SunriseAndSunset:
		m.IsLoading = 0
		m.SunriseAndSunset = msg.SunriseAndSunsetResponse
		return m, DoTick()

	case RefreshWeatherMsg:
		if m.IsLoading != 0 {
			return m, DoWeatherRefreshTick()
		}

		return m, tea.Batch(
			FetchWeather(m.Location, m.Config.Units),
			DoWeatherRefreshTick(),
		)

	case TickMsg:
		m.CurrentTime = time.Time(msg)
		return m, DoTick()

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if msg.String() == "r" {
			if m.Err != nil {
				m.IsLoading = 2
			}
			m.Err = nil
			return m, FetchWeather(m.Location, m.Config.Units)
		}

		if msg.String() == "s" {
			m.Config.HideSunBar = !m.Config.HideSunBar
			cfg.SaveConfig(m.Config)
			return m, DoTick()
		}

	case ErrMsg:
		m.IsLoading = -1
		var netErr net.Error
		if errors.As(msg.Err, &netErr) {
			if netErr.Timeout() {
				m.Err = fmt.Errorf("Network timeout.")
			} else {
				m.Err = fmt.Errorf("You're offline.")
			}
		} else {
			m.Err = msg.Err
		}

		return m, nil
	}

	return m, nil
}

func (m Model) View() string {

	if m.Err != nil || m.IsLoading == -1 {
		header := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, m.Err.Error())
		version := m.renderVersion()

		footer := lipgloss.Place(
			m.Width,
			1,
			lipgloss.Right,
			lipgloss.Bottom,
			version,
		)
		return lipgloss.JoinVertical(lipgloss.Left, header, footer)
	}

	if m.IsLoading == 1 {
		header := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, "📍 Loading location info, please wait.")
		version := m.renderVersion()

		footer := lipgloss.Place(
			m.Width,
			1,
			lipgloss.Right,
			lipgloss.Bottom,
			version,
		)
		return lipgloss.JoinVertical(lipgloss.Left, header, footer)
	}

	if m.IsLoading == 2 {
		header := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, "⛅ Loading weather info, please wait.")

		version := m.renderVersion()

		footer := lipgloss.Place(
			m.Width,
			1,
			lipgloss.Right,
			lipgloss.Bottom,
			version,
		)
		return lipgloss.JoinVertical(lipgloss.Left, header, footer)
	}

	if m.IsLoading == 3 {
		header := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, "🌇 Loading sunrise and sunset data, please wait.")

		version := m.renderVersion()

		footer := lipgloss.Place(
			m.Width,
			1,
			lipgloss.Right,
			lipgloss.Bottom,
			version,
		)
		return lipgloss.JoinVertical(lipgloss.Left, header, footer)
	}

	weatherArt := ascii.GetCurrentWeatherArt(m.Weather.CurrentWeather.WeatherCode)

	var cleanArtLines []string
	for _, line := range strings.Split(strings.TrimSpace(weatherArt), "\n") {
		cleanArtLines = append(cleanArtLines, strings.TrimSpace(line))
	}
	cleanArt := strings.Join(cleanArtLines, "\n")

	timeStr := m.CurrentTime.Format("15:04:05")
	dateStr := fmt.Sprintf(
		"%s, %s %d",
		m.CurrentTime.Weekday(),
		m.CurrentTime.Month(),
		m.CurrentTime.Day(),
	)

	var loc string

	if len(m.Location.Region) == 0 || m.Location.Region == "" {
		loc = fmt.Sprintf("📍 %s", m.Location.City)
	} else if len(m.Location.City) == 0 || m.Location.City == "" {
		loc = fmt.Sprintf("📍 %s", m.Location.Region)
	} else {
		loc = fmt.Sprintf("📍 %s, %s", m.Location.Region, m.Location.City)
	}

	weatherName := weather.GetCurrentWeatherName(m.Weather.CurrentWeather.WeatherCode)

	temp := m.Config.FormatTemp(m.Weather.CurrentWeather.Temperature2m)
	var sunBar string
	// TODO: Replacing sunbar with hourly/daily weather or graph
	if !m.Config.HideSunBar {
		sunBar = m.renderSunBar(m.Width)
	}

	mainContent := lipgloss.JoinVertical(lipgloss.Center, loc, "", cleanArt, "", temp, weatherName, "", timeStr, dateStr)

	footer := lipgloss.JoinVertical(lipgloss.Center, sunBar)
	footerHeight := lipgloss.Height(footer)

	centeredMain := lipgloss.Place(
		m.Width,
		m.Height-footerHeight,
		lipgloss.Center,
		lipgloss.Center,
		mainContent,
	)

	return lipgloss.JoinVertical(lipgloss.Left, centeredMain, footer)
}

func (m Model) renderVersion() string {
	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1)
	return versionStyle.Render(m.Version)
}

func (m Model) renderSunBar(width int) string {
	progress := m.getSunProgress()

	barWidth := width / 6
	if barWidth < 20 {
		barWidth = 20
	}

	filledWidth := int(progress * float64(barWidth))
	if filledWidth > barWidth {
		filledWidth = barWidth
	}
	if filledWidth < 0 {
		filledWidth = 0
	}

	bar := strings.Repeat("▓", filledWidth) + strings.Repeat("░", barWidth-filledWidth)

	st := m.SunriseAndSunset.Daily.Sunrise[0]
	et := m.SunriseAndSunset.Daily.Sunset[0]
	sunriseTime := st[strings.Index(st, "T")+1:]
	sunsetTime := et[strings.Index(et, "T")+1:]

	content := fmt.Sprintf("%s  %s  %s", sunriseTime, bar, sunsetTime)

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(content)
}

func (m Model) getSunProgress() float64 {
	layout := "2006-01-02T15:04"
	sunrise, _ := time.Parse(layout, m.SunriseAndSunset.Daily.Sunrise[0])
	sunset, _ := time.Parse(layout, m.SunriseAndSunset.Daily.Sunset[0])

	if m.CurrentTime.Before(sunrise) {
		return 0.0
	}
	if m.CurrentTime.After(sunset) {
		return 1.0
	}

	totalDuration := sunset.Sub(sunrise).Seconds()
	elapsedDuration := m.CurrentTime.Sub(sunrise).Seconds()

	return elapsedDuration / totalDuration
}
