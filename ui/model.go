package ui

import (
	"fmt"
	"time"

	"github.com/zenpaw-labs/skypaw/ascii"
	"github.com/zenpaw-labs/skypaw/network/geocoding"
	"github.com/zenpaw-labs/skypaw/network/weather"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	// Weather
	City     string
	Weather  weather.WeatherResponse
	Location geocoding.LocationInfo

	// Status
	CurrentTime time.Time
	IsLoading   int
	Err         error

	// Window
	Width  int
	Height int
}

func InitialModel(city string) Model {
	return Model{
		City:        city,
		CurrentTime: time.Now(),
		IsLoading:   1,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		FetchWeatherWithAutoLocation(),
		DoTick(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

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

	case WeatherMsg:
		m.Weather = msg.Data
		m.Location = msg.LocationInfo
		m.IsLoading = 0
		return m, DoTick()

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
		return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, "⏳ Loading weather info . . .")
	}

	if m.Err != nil {
		return "❌ Error: " + m.Err.Error()
	}

	if m.IsLoading == 1 {
		return "\n  ⏳ Loading weather info . . ."
	}

	weatherArt := ascii.GetCurrentWeatherArt(m.Weather.CurrentWeather.WeatherCode)

	timeStr := m.CurrentTime.Format("15:04:05")

	s := fmt.Sprintf(
		"\n  🏙  City: %s\n\n"+
			"%s\n\n"+
			"  🌡  Temperature: %.1f°C\n"+
			"  🕒 Time: %s\n\n"+
			"  ('q' for exit)",
		m.Location.Name,
		weatherArt,
		m.Weather.CurrentWeather.Temperature2m,
		timeStr,
	)

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		s,
	)
}
