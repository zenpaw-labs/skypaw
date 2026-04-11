package ui

import (
	"skypaw/network/weather"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func FetchWeather(city string) tea.Cmd {
	return func() tea.Msg {
		res, info, err := weather.GetCurrentWeather(city)
		if err != nil {
			return ErrMsg{err}
		}
		return WeatherMsg{Data: res, LocationInfo: info}
	}
}

func DoTick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
