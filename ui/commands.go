package ui

import (
	"time"

	"github.com/zenpaw-labs/skypaw/network/geocoding"
	"github.com/zenpaw-labs/skypaw/network/weather"
	"github.com/zenpaw-labs/skypaw/utils/location_utils"

	tea "github.com/charmbracelet/bubbletea"
)

func FetchWeather(location geocoding.LocationInfo) tea.Cmd {
	return func() tea.Msg {
		res, info, err := weather.GetCurrentWeather(location.Name)
		if err != nil {
			return ErrMsg{err}
		}
		return WeatherMsg{Data: res, LocationInfo: info}
	}
}

func FetchLocationByName(location string) tea.Cmd {
	return func() tea.Msg {
		return GeocodingMsg{Data: geocoding.SearchLocation(location)}
	}
}

func FetchLocation() tea.Cmd {
	return func() tea.Msg {
		l, err := location_utils.GetLocationFromOs()
		if err != nil {
			return ErrMsg{err}
		}
		return GeocodingMsg{l}
	}
}

func DoTick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
