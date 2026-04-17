package ui

import (
	"time"

	"github.com/zenpaw-labs/skypaw/network/geocoding"
	"github.com/zenpaw-labs/skypaw/network/weather"
)

type WeatherMsg struct {
	Data         weather.WeatherResponse
	LocationInfo geocoding.LocationInfo
}

type ErrMsg struct {
	Err error
}

type TickMsg time.Time
