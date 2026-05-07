package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zenpaw-labs/skypaw/ascii"
	"github.com/zenpaw-labs/skypaw/network/geocoding"
	"github.com/zenpaw-labs/skypaw/network/weather"
)

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
	optionalProvider *int
	customCity       string

	// Other
	Version string
}

func InitialModel(optionalProvider *int, version string, city string) Model {
	return Model{
		customCity:       city,
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
		m.IsLoading = 3
		return m, FetchSunriseAndSunset(m.Location)

	case SunriseAndSunset:
		m.IsLoading = 0
		m.SunriseAndSunset = msg.SunriseAndSunsetResponse
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

	if m.Err != nil {
		return "❌ Error: " + m.Err.Error()
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
	
	loc := fmt.Sprintf("📍 %s, %s", m.Location.Admin1, m.Location.Name)
	weatherName := weather.GetCurrentWeatherName(m.Weather.CurrentWeather.WeatherCode)
	temp := fmt.Sprintf("%.1f°C", m.Weather.CurrentWeather.Temperature2m)
	sunBar := m.renderSunBar(m.Width)

	mainContent := lipgloss.JoinVertical(
        lipgloss.Center,
        loc,
        "",
        cleanArt,
        "",
        temp,
        weatherName,
        "",
        timeStr,
        dateStr,
    )

    centeredMain := lipgloss.Place(
        m.Width,
        m.Height-3, 
        lipgloss.Center,
        lipgloss.Center,
        mainContent,
    )


    footerContent := lipgloss.NewStyle().
	Width(m.Width).Align(lipgloss.Center).
	Render(lipgloss.JoinVertical(
		lipgloss.Center,
		sunBar,
	))

    return lipgloss.JoinVertical(lipgloss.Left, centeredMain, footerContent)
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
    
    filledWidth := int(progress * float64(barWidth))
    if filledWidth > barWidth { filledWidth = barWidth }
    
    filledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
    emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#3C3C3C"))
    
    bar := filledStyle.Render(strings.Repeat("━", filledWidth)) + 
    emptyStyle.Render(strings.Repeat("━", barWidth-filledWidth))
          
    st := m.SunriseAndSunset.Daily.Sunrise[0]
    et := m.SunriseAndSunset.Daily.Sunset[0]
    sunriseTime := st[strings.Index(st, "T")+1:]
    sunsetTime := et[strings.Index(et, "T")+1:]

    labels := fmt.Sprintf("🌅 %s %s 🌇 %s", sunriseTime, bar, sunsetTime)
    return labels
}

func (m Model) getSunProgress() float64 {
    layout := "2006-01-02T15:04"
    sunrise, _ := time.Parse(layout, m.SunriseAndSunset.Daily.Sunrise[0])
    sunset, _ := time.Parse(layout, m.SunriseAndSunset.Daily.Sunset[0])

    if m.CurrentTime.Before(sunrise) { return 0.0 }
    if m.CurrentTime.After(sunset) { return 1.0 }

    totalDuration := sunset.Sub(sunrise).Seconds()
    elapsedDuration := m.CurrentTime.Sub(sunrise).Seconds()

    return elapsedDuration / totalDuration
}