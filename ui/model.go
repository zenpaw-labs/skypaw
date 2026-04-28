package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zenpaw-labs/skypaw/ascii"
	"github.com/zenpaw-labs/skypaw/network/geocoding"
	"github.com/zenpaw-labs/skypaw/network/weather"
)

type Model struct {
	// Weather
	City     string
	Weather  weather.WeatherResponse
	Location geocoding.LocationInfo

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
	optionalProvider *int
	customCity string

	// Other
	Version          string
}

func InitialModel(optionalProvider *int, version string, city string) Model {
	return Model{
		customCity: city,
		optionalProvider: optionalProvider,
		Version:          version,
		CurrentTime:      time.Now(),
		IsLoading:        1,
	}
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	if m.customCity != "" {
		cmds = append(cmds, FetchLocationByName(m.customCity))
	} else {
		cmds = append(cmds, FetchLocation(m.optionalProvider))
	}
	cmds = append(cmds, DoTick())
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case GeocodingMsg:
		m.Location = msg.Data
		m.IsLoading = 2
		return m, FetchWeather(m.Location)
	case WeatherMsg:
		m.Weather = msg.Data
		m.Location = msg.LocationInfo
		m.IsLoading = 0
		return m, DoTick()

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case TickMsg:
		m.CurrentTime = time.Time(msg)
		return m, DoTick()

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case ErrMsg:
		m.Err = msg.Err
		m.IsLoading = -1
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {

	if m.Err != nil {
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, "❌ Error: "+m.Err.Error())
	}

	if m.IsLoading == 1 {
		header := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, "📍 Loading location info, please wait.")
		versionStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 2)

		versionBlock := versionStyle.Render(m.Version)
		footer := lipgloss.Place(
			m.Width,
			1,
			lipgloss.Right,
			lipgloss.Bottom,
			versionBlock,
		)
		return lipgloss.JoinVertical(lipgloss.Left, header, footer)
	}

	if m.IsLoading == 2 {
		header := lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, "⛅ Loading weather info, please wait.")

		versionStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFFFFF")).
			Bold(true).
			Padding(0, 2)

		versionBlock := versionStyle.Render(m.Version)
		footer := lipgloss.Place(
			m.Width,
			1,
			lipgloss.Right,
			lipgloss.Bottom,
			versionBlock,
		)
		return lipgloss.JoinVertical(lipgloss.Left, header, footer)
	}

	if m.Err != nil {
		return "❌ Error: " + m.Err.Error()
	}

	weatherArt := ascii.GetCurrentWeatherArt(m.Weather.CurrentWeather.WeatherCode)

	timeStr := m.CurrentTime.Format("15:04:05")
	dateStr := fmt.Sprintf(
		"%s, %s %d",
		m.CurrentTime.Weekday(),
		m.CurrentTime.Month(),
		m.CurrentTime.Day(),
	)
	loc := fmt.Sprintf("📍 %s, %s", m.Location.Admin1, m.Location.Name)
	s := fmt.Sprintf(
		"%s\n\n"+ // Location data
			"%s\n\n"+ // ASCII Art
			"%.1f°C\n\n"+ // Weather Temperature
			"%s\n"+ // Time
			"%s\n", // Weekday
		loc,
		weatherArt,
		m.Weather.CurrentWeather.Temperature2m,
		timeStr,
		dateStr,
	)

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		s,
	)
}
