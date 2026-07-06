package ui

import (
	"time"

	"github.com/zenpaw-labs/skypaw/network/geocoding"
	"github.com/zenpaw-labs/skypaw/network/weather"
	"github.com/zenpaw-labs/skypaw/utils/cfg"
	"github.com/zenpaw-labs/skypaw/utils/location_utils"

	tea "github.com/charmbracelet/bubbletea"
)

func FetchWeather(location geocoding.LocationInfo, units cfg.UnitSystem) tea.Cmd {
	return func() tea.Msg {
		res, info, err := weather.GetCurrentWeatherByLocationInfo(location, units)
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

func FetchLocation(optionalProvider *int) tea.Cmd {
	return func() tea.Msg {
		l, err := location_utils.GetLocation(optionalProvider)
		if err != nil {
			return ErrMsg{err}
		}
		return GeocodingMsg{l}
	}
}

func FetchSunriseAndSunset(l geocoding.LocationInfo) tea.Cmd {
	return func() tea.Msg {
		s, err := weather.GetSunriseAndSunset(l)
		if err != nil {
			return ErrMsg{err}
		}
		return SunriseAndSunset{s}
	}
}

func DoWeatherRefreshTick() tea.Cmd {
	return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
		return RefreshWeatherMsg(t)
	})
}

func DoTick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
