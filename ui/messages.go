package ui

import (
	"skypaw/network/geocoding"
	"skypaw/network/weather"
	"time"
)

type WeatherMsg struct {
	Data         weather.WeatherResponse
	LocationInfo geocoding.LocationInfo
}

type ErrMsg struct {
	Err error
}

type TickMsg time.Time
