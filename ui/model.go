package ui

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zenpaw-labs/skypaw/network/geocoding"
	"github.com/zenpaw-labs/skypaw/network/weather"
	"github.com/zenpaw-labs/skypaw/utils/cfg"
)

//TODO: Interactive location picker

const (
	LoadingError     = -1
	LoadingCompleted = 0
	LoadingLocation  = 1
	LoadingWeather   = 2
	LoadingSunrise   = 3
	LoadingHourly    = 4
)

type Model struct {
	// Weather
	City             string
	Weather          weather.WeatherResponse
	SunriseAndSunset weather.SunriseAndSunsetResponse
	Hourly           weather.HourlyWeatherResponse
	Location         geocoding.LocationInfo

	// Status
	CurrentTime    time.Time
	CurrentWeekday time.Weekday
	CurrentMonth   time.Month
	IsLoading      int
	Err            error
	spinner        spinner.Model

	// Window
	Width  int
	Height int

	// User config
	Config cfg.UserConfig

	// Other
	Version string
}

func InitialModel(cfg cfg.UserConfig, version string) Model {
	s := spinner.New()
	s.Spinner = spinner.MiniDot

	return Model{
		Config:      cfg,
		Version:     version,
		CurrentTime: time.Now(),
		IsLoading:   LoadingLocation,
		spinner:     s,
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
	cmds = append(cmds, m.spinner.Tick)
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case GeocodingMsg:
		m.Location = msg.Data
		if m.IsLoading != LoadingCompleted {
			m.IsLoading = LoadingWeather
		}
		return m, FetchWeather(m.Location, m.Config.Units)

	case WeatherMsg:
		m.Weather = msg.Data
		m.Location = msg.LocationInfo
		if m.IsLoading != LoadingCompleted {
			m.IsLoading = LoadingHourly
		}
		return m, FetchHourlyWeather(m.Location)

	case HourlyWeatherMsg:
		m.Hourly = msg.HourlyWeatherResponse
		if m.IsLoading != LoadingCompleted {
			m.IsLoading = LoadingSunrise
		}
		return m, FetchSunriseAndSunset(m.Location)

	case SunriseAndSunsetMsg:
		m.IsLoading = LoadingCompleted
		m.SunriseAndSunset = msg.SunriseAndSunsetResponse
		return m, DoTick()

	case RefreshWeatherMsg:
		if m.IsLoading != LoadingCompleted {
			return m, DoWeatherRefreshTick()
		}

		return m, tea.Batch(
			FetchWeather(m.Location, m.Config.Units),
			FetchHourlyWeather(m.Location),
			DoWeatherRefreshTick(),
		)

	case TickMsg:
		m.CurrentTime = time.Time(msg)
		return m, DoTick()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if msg.String() == "r" {
			m.Err = nil

			if m.Location.City == "" && m.Location.Region == "" {
				m.IsLoading = LoadingLocation
				if m.Config.UserCity != "" {
					return m, FetchLocationByName(m.Config.UserCity)
				}
				return m, FetchLocation(m.Config, m.Config.WindowsLocalLocationDetection)
			}
			m.IsLoading = LoadingWeather
			return m, FetchWeather(m.Location, m.Config.Units)
		}

		if msg.String() == "s" {
			m.Config.HideSunBar = !m.Config.HideSunBar
			cfg.SaveConfig(m.Config)
			return m, DoTick()
		}

	case ErrMsg:
		m.IsLoading = LoadingError
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

	loadingText := map[int]string{
		LoadingLocation: "📍 Detecting location",
		LoadingWeather:  "⛅ Fetching weather",
		LoadingHourly:   "🕐 Loading hourly data",
		LoadingSunrise:  "🌇 Loading sunrise data",
	}

	if m.IsLoading == LoadingError || m.Err != nil {
		content := fmt.Sprintf("%s", m.Err.Error())
		version := m.renderVersion()
		header := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
		footer := lipgloss.Place(m.Width, 1, lipgloss.Right, lipgloss.Bottom, version)
		return lipgloss.JoinVertical(lipgloss.Left, header, footer)
	}

	if m.IsLoading != LoadingCompleted && m.IsLoading != LoadingError {
		text := loadingText[m.IsLoading]
		content := fmt.Sprintf("%s %s", m.spinner.View(), text)

		version := m.renderVersion()
		header := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
		footer := lipgloss.Place(m.Width, 1, lipgloss.Right, lipgloss.Bottom, version)
		return lipgloss.JoinVertical(lipgloss.Left, header, footer)
	}

	weatherArt := GetCurrentWeatherArt(m.Weather.CurrentWeather.WeatherCode)

	rawArt := strings.Trim(weatherArt.Art, "\n\r")
	lines := strings.Split(rawArt, "\n")

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

	var croppedLines []string
	artWidth := 0
	for _, line := range lines {
		line = strings.TrimRight(line, "\r ")
		if len(line) >= minSpaces {
			cropped := line[minSpaces:]
			croppedLines = append(croppedLines, cropped)

			w := lipgloss.Width(cropped)
			if w > artWidth {
				artWidth = w
			}
		} else {
			croppedLines = append(croppedLines, "")
		}
	}

	for i, line := range croppedLines {
		w := lipgloss.Width(line)
		if w < artWidth {
			croppedLines[i] = line + strings.Repeat(" ", artWidth-w)
		}
	}

	cleanRectArt := strings.Join(croppedLines, "\n")

	artStyle := lipgloss.NewStyle().
		Width(artWidth).
		Align(lipgloss.Left)

	var art string
	if m.Config.ColorfulTUI {
		art = artStyle.Foreground(weatherArt.Color).Render(cleanRectArt)
	} else {
		art = artStyle.Render(cleanRectArt)
	}
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

	mainContent := lipgloss.JoinVertical(lipgloss.Center, loc, "", art, "", temp, weatherName, "", timeStr, dateStr)

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

func (m Model) renderHourlyTemperature() string {
	return ""
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
	if len(m.SunriseAndSunset.Daily.Sunrise) == 0 ||
		len(m.SunriseAndSunset.Daily.Sunset) == 0 {
		return ""
	}
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

	content := fmt.Sprintf("🔆 %s  %s  %s 🌒", sunriseTime, bar, sunsetTime)

	return lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(content)
}

func (m Model) getSunProgress() float64 {
	if len(m.SunriseAndSunset.Daily.Sunrise) == 0 ||
		len(m.SunriseAndSunset.Daily.Sunset) == 0 {
		return 0.0
	}
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
