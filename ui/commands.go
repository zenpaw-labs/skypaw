package ui

import (
	"time"

	"github.com/zenpaw-labs/skypaw/network/geocoding"
	"github.com/zenpaw-labs/skypaw/network/weather"
	"github.com/zenpaw-labs/skypaw/utils"
	"github.com/zenpaw-labs/skypaw/utils/location_utils"

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

func FetchWeatherWithAutoLocation() tea.Cmd {
	return func() tea.Msg {
		c, err := location_utils.GetLocationFromOs(utils.GetRuntimeOs())
		if err != nil {
			return ErrMsg{err}
		}

		w, _, err := weather.GetCurrentWeatherFromCoordinates(c)
		if err != nil {
			return ErrMsg{err}
		}
		location, err := geocoding.GetLocationFromCoords(c)
		if err != nil {
			return ErrMsg{err}
		}
		return WeatherMsg{Data: w, LocationInfo: location}
	}
}

func DoTick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
