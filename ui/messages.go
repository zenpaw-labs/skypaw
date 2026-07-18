package ui

import (
	"time"

	"github.com/zenpaw-labs/skypaw/network/geocoding"
	"github.com/zenpaw-labs/skypaw/network/weather"
)

type GeocodingMsg struct {
	Data geocoding.LocationInfo
}

type WeatherMsg struct {
	Data         weather.WeatherResponse
	LocationInfo geocoding.LocationInfo
}

type SunriseAndSunsetMsg struct {
	weather.SunriseAndSunsetResponse
}

type HourlyWeatherMsg struct {
	weather.HourlyWeatherResponse
}

type ErrMsg struct {
	Err error
}

type TickMsg time.Time

type RefreshWeatherMsg time.Time
