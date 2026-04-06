package weather

import (
	"fmt"
	"net/url"
	"skypaw/network"
	"skypaw/network/geocoding"
	"strconv"
)

type WeatherInfo struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	Elevation            float64 `json:"elevation"`
	GenerationTimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Current              any     `json:"current"`
	Hourly               any     `json:"hourly"`
	HourlyUnits          any     `json:"hourly_units"`
	Daily                any     `json:"daily"`
	DailyUnits           any     `json:"daily_units"`
}

func SearchWeather(locationInfo geocoding.LocationInfo) WeatherInfo {
	var (
		wInfo WeatherInfo
	)

	values := url.Values{}
	values.Add("latitude", strconv.FormatFloat(locationInfo.Latitude, 'f', -1, 64))
	values.Add("longitude", strconv.FormatFloat(locationInfo.Longitude, 'f', -1, 64))
	fullUrl := network.WeatherEndpointApi + "forecast?" + values.Encode()
	fmt.Println(fullUrl)
	return wInfo
}
