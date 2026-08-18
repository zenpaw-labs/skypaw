package ui

import (
	"errors"
	"fmt"
	"math"
	"net"
	"strings"
	"time"

	"charm.land/log/v2"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
	"github.com/zenpaw-labs/skypaw/network/geocoding"
	"github.com/zenpaw-labs/skypaw/network/weather"
	"github.com/zenpaw-labs/skypaw/utils/cfg"
)

//TODO: Interactive location picker

const (
	// Loading
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
	CurrentTime     time.Time
	CurrentWeekday  time.Weekday
	CurrentMonth    time.Month
	IsLoading       int
	isOfflineMode   bool
	offlineModeData weather.WeatherCache
	Err             error
	spinner         spinner.Model

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
		return m, FetchHourlyWeather(m.Location, m.Config.Units)

	case HourlyWeatherMsg:
		m.Hourly = msg.HourlyWeatherResponse

		if err := convertToLocalTime(&m.Hourly); err != nil {
			log.Error("Failed to convert hourly time", "err", err)
		}

		if m.IsLoading != LoadingCompleted {
			m.IsLoading = LoadingSunrise
		}
		return m, FetchSunriseAndSunset(m.Location)

	case SunriseAndSunsetMsg:
		m.IsLoading = LoadingCompleted
		m.SunriseAndSunset = msg.SunriseAndSunsetResponse

		if m.Config.UseWeatherCache {
			go weather.Save(m.Location, m.Weather, m.Hourly, m.SunriseAndSunset, m.Config.Units)
		}
		return m, DoTick()

	case RefreshWeatherMsg:
		if m.IsLoading != LoadingCompleted {
			return m, DoWeatherRefreshTick()
		}

		return m, tea.Batch(
			FetchWeather(m.Location, m.Config.Units),
			FetchHourlyWeather(m.Location, m.Config.Units),
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
			m.isOfflineMode = false
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
			m.Config.HideDiagram = !m.Config.HideDiagram
			cfg.SaveConfig(m.Config)
			return m, DoTick()
		}

		if msg.String() == "h" {
			if m.Config.ShowHints {
				m.Config.ShowHints = false
			} else {
				m.Config.ShowHints = true
			}
			cfg.SaveConfig(m.Config)
		}

	case ErrMsg:
		m.IsLoading = LoadingError
		var netErr net.Error

		if errors.As(msg.Err, &netErr) {
			if m.tryRestoreFromCache() {
				return m, nil
			}

			t := "You're offline."
			if netErr.Timeout() {
				t = "Network timeout."
			}

			if m.Config.ColorfulTUI {
				t = lipgloss.NewStyle().Foreground(ColorOfflineOrTimeout).Render(t)
			}
			m.Err = fmt.Errorf("%s", t)
		} else {
			m.Err = msg.Err
		}

		log.Error(m.Err)
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
		var content string
		if m.Config.ColorfulTUI {
			s := lipgloss.NewStyle().Foreground(ColorError)
			content = fmt.Sprintf("%s", s.Render(m.Err.Error()))
		} else {
			content = fmt.Sprintf("%s", m.Err.Error())
		}
		version := m.renderVersion()
		header := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, content)
		footer := lipgloss.Place(m.Width, 1, lipgloss.Right, lipgloss.Bottom, version)
		return lipgloss.JoinVertical(lipgloss.Left, header, footer)
	}

	if m.IsLoading != LoadingCompleted && m.IsLoading != LoadingError {
		text := loadingText[m.IsLoading]
		if m.Config.ColorfulTUI {
			textStyle := lipgloss.NewStyle().Foreground(ColorLoading)
			text = textStyle.Render(text)
		}

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

	var temp string
	if m.Config.Units == cfg.Imperial {
		temp = fmt.Sprintf("%.1f°F", m.Weather.CurrentWeather.Temperature2m)
	} else {
		temp = fmt.Sprintf("%.1f°C", m.Weather.CurrentWeather.Temperature2m)
	}
	var diagram string

	if !m.Config.HideDiagram {
		diagram = m.renderHourlyTemperature()
	}

	var art string
	artStyle := lipgloss.NewStyle().
		Width(artWidth).
		Align(lipgloss.Left)
	if m.Config.ColorfulTUI {

		locStyle := lipgloss.NewStyle().Foreground(ColorLocation)
		tempStyle := lipgloss.NewStyle().Foreground(ColorTemp)
		weatherNameStyle := lipgloss.NewStyle().Foreground(ColorWeatherName)
		timeStyle := lipgloss.NewStyle().Foreground(ColorTime)
		dateStyle := lipgloss.NewStyle().Foreground(ColorDate)

		loc = locStyle.Render(loc)
		temp = tempStyle.Render(temp)
		weatherName = weatherNameStyle.Render(weatherName)
		timeStr = timeStyle.Render(timeStr)
		dateStr = dateStyle.Render(dateStr)

		art = artStyle.Foreground(weatherArt.Color).Render(cleanRectArt)
	} else {
		art = artStyle.Render(cleanRectArt)
	}

	var hints string
	if m.Config.ShowHints {
		hints = "\nPress 'Q' or Ctrl+C for exit.\nTo show / hide diagram press 'S'.\nTo hide this message press 'H'."
		hintsStyle := lipgloss.NewStyle().Foreground(ColorHints)
		hints = hintsStyle.Render(hints)

	}

	var offlineModeNotice string
	if m.isOfflineMode {
		offlineModeNotice = fmt.Sprintf("\nOffline mode • Last update at %s", m.offlineModeData.LastUpdatedAt.Format("15:04 (02.01.2006)"))
		style := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(m.Width)
		if m.Config.ColorfulTUI {
			style = style.Foreground(ColorOfflineOrTimeout)
		}
		offlineModeNotice = style.Render(offlineModeNotice) + "\n"
	}
	mainContent := lipgloss.JoinVertical(lipgloss.Center, loc, "", art, "", temp, weatherName, "", timeStr, dateStr, offlineModeNotice, hints)

	footer := lipgloss.JoinVertical(lipgloss.Center, diagram)
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
	if len(m.Hourly.Hourly.Temperature2m) == 0 {
		return ""
	}

	slicedTemps, slicedTimes := getDiagramSlices(
		m.Hourly.Hourly.Temperature2m,
		m.Hourly.Hourly.Time,
		getLocalHourIndex(m.Hourly.Hourly.Time),
		m.Config.DiagramHoursBefore,
		m.Config.DiagramHoursAfter,
	)

	if len(slicedTemps) == 0 {
		return ""
	}

	height := 3

	width := m.Width / 3
	if width > 30 {
		width = 30
	}
	if width < 10 {
		width = 10
	}

	xAxisFormatter := func(v float64) string {
		idx := int(math.Round(v))
		if idx >= 0 && idx < len(slicedTimes) {
			if t, err := time.Parse("2006-01-02T15:04", slicedTimes[idx]); err == nil {
				return t.Format("15:04")
			}
		}
		return ""
	}
	graph := asciigraph.Plot(
		slicedTemps,
		asciigraph.Width(width),
		asciigraph.Height(height),
		asciigraph.XAxisRange(0, float64(len(slicedTemps)-1)),
		asciigraph.XAxisValueFormatter(xAxisFormatter),
		asciigraph.Precision(0),
	)

	return lipgloss.PlaceHorizontal(
		m.Width,
		lipgloss.Center,
		padLinesToEqualWidth(graph),
	)
}

func (m Model) renderVersion() string {
	versionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 1)
	return versionStyle.Render(m.Version)
}

func convertToLocalTime(resp *weather.HourlyWeatherResponse) error {
	layout := "2006-01-02T15:04"
	for i, timeStr := range resp.Hourly.Time {
		utcTime, err := time.ParseInLocation(layout, timeStr, time.UTC)
		if err != nil {
			return fmt.Errorf("failed to parse time %s at index %d: %w", timeStr, i, err)
		}

		resp.Hourly.Time[i] = utcTime.Local().Format(layout)
	}
	return nil
}

func getLocalHourIndex(times []string) int {
	now := time.Now().Truncate(time.Hour)

	layout := "2006-01-02T15:00"

	target := now.Format(layout)

	for i, t := range times {
		if t == target {
			return i
		}
	}

	return -1
}

func getDiagramSlices(
	temps []float64,
	times []string,
	hourIndex int,
	hoursBefore int,
	hoursAfter int,
) ([]float64, []string) {
	totalLen := len(temps)

	if totalLen == 0 || hourIndex < 0 || hourIndex >= totalLen {
		return nil, nil
	}

	start := hourIndex - hoursBefore
	if start < 0 {
		start = 0
	}

	end := hourIndex + hoursAfter + 1
	if end > totalLen {
		end = totalLen
	}

	if start >= end {
		return nil, nil
	}

	slicedTemps := make([]float64, end-start)
	copy(slicedTemps, temps[start:end])

	slicedTimes := make([]string, end-start)
	copy(slicedTimes, times[start:end])

	return slicedTemps, slicedTimes
}

func padLinesToEqualWidth(s string) string {
	lines := strings.Split(s, "\n")
	maxWidth := 0

	for _, line := range lines {
		if w := lipgloss.Width(line); w > maxWidth {
			maxWidth = w
		}
	}

	for i, line := range lines {
		w := lipgloss.Width(line)
		if w < maxWidth {
			lines[i] = line + strings.Repeat(" ", maxWidth-w)
		}
	}

	return strings.Join(lines, "\n")
}

func (m *Model) tryRestoreFromCache() bool {
	if !m.Config.UseWeatherCache {
		return false
	}

	data, err := weather.Load(m.Config)
	if err != nil {
		log.Error("Failed to load weather cache", "err", err)
		return false
	}
	m.offlineModeData = *data
	m.Location = m.offlineModeData.Location
	m.Weather = m.offlineModeData.Weather
	m.Hourly = m.offlineModeData.Hourly
	m.SunriseAndSunset = m.offlineModeData.SunriseSunset

	m.IsLoading = LoadingCompleted
	m.isOfflineMode = true
	m.Err = nil
	return true
}
