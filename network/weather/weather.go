package weather

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zenpaw-labs/skypaw/network"
	"github.com/zenpaw-labs/skypaw/network/geocoding"
	"github.com/zenpaw-labs/skypaw/utils/cfg"
)

var (
	weatherCodes = map[int]string{
		0:  "Clear Sky",
		1:  "Mainly Clear",
		2:  "Partly Cloudy",
		3:  "Overcast",
		45: "Fog",
		48: "Depositing rime fog",
		51: "Drizzle: Light intensity",
		53: "Drizzle: Moderate intensity",
		55: "Drizzle: Dense intensity ",
		56: "Freezing Drizzle: Light intensity",
		57: "Freezing Drizzle: Dense intensity",
		61: "Rain: Slight intensity",
		63: "Rain: Moderate intensity",
		65: "Rain: Heavy intensity",
		66: "Freezing Rain: Light intensity",
		67: "Freezing Rain: Heavy intensity",
		71: "Snowfall: Slight intensity",
		73: "Snowfall: Moderate intensity",
		75: "Snowfall: Heavy intensity",
		77: "Snow Grains",
		80: "Rain Shower: Slight intensity",
		81: "Rain Shower: Moderate intensity",
		82: "Rain Shower: Violent intensity",
		85: "Snow Shower: Slight intensity",
		86: "Snow Shower: Heavy intensity",
		95: "Thunderstorm: Slight or moderate",
		96: "Thunderstorm with slight hail",
		99: "Thunderstorm with heavy hail",
	}
)

type WeatherResponse struct {
	Latitude             float64        `json:"latitude"`
	Longitude            float64        `json:"longitude"`
	GenerationTimeMs     float64        `json:"generationtime_ms"`
	Timezone             string         `json:"timezone"`
	TimezoneAbbreviation string         `json:"timezone_abbreviation"`
	Elevation            float64        `json:"elevation"`
	CurrentWeather       CurrentWeather `json:"current"`
}

type CurrentWeather struct {
	Time          string  `json:"time"`
	Interval      int     `json:"interval"`
	Temperature2m float64 `json:"temperature_2m"`
	WeatherCode   int     `json:"weather_code"`
	WindSpeed10m  float64 `json:"wind_speed_10m"`
	IsDay         int     `json:"is_day"`
}

type HourlyWeatherResponse struct {
	Hourly struct {
		Time          []string  `json:"time"`
		Temperature2m []float64 `json:"temperature_2m"`
	} `json:"hourly"`
}

type SunriseAndSunsetResponse struct {
	/*
		Struct generated for according to response from OpenMeteo daily weather
		Learn more: https://open-meteo.com/en/docs#daily_weather_variables
	*/
	Latitude             float64    `json:"latitude"`
	Longitude            float64    `json:"longitude"`
	GenerationtimeMS     float64    `json:"generationtime_ms"`
	UTCOffsetSeconds     int64      `json:"utc_offset_seconds"`
	Timezone             string     `json:"timezone"`
	TimezoneAbbreviation string     `json:"timezone_abbreviation"`
	Elevation            int64      `json:"elevation"`
	DailyUnits           DailyUnits `json:"daily_units"`
	Daily                Daily      `json:"daily"`
}

type Daily struct {
	Time    []string `json:"time"`
	Sunrise []string `json:"sunrise"`
	Sunset  []string `json:"sunset"`
}

type DailyUnits struct {
	Time    string `json:"time"`
	Sunrise string `json:"sunrise"`
	Sunset  string `json:"sunset"`
}

func GetCurrentWeatherByLocationInfo(locationInfo geocoding.LocationInfo, unitSystem cfg.UnitSystem) (WeatherResponse, geocoding.LocationInfo, error) {
	/*
		Requests generated for according to API Scheme of current weather by OpenMeteo
		Docs of current weather API: https://open-meteo.com/en/docs#current_weather
	*/
	var (
		weatherResponse WeatherResponse
	)
	rq := []string{"temperature_2m", "is_day", "weather_code", "wind_speed_10m"}
	// TODO: Pressure, Wind (with arrow of direction) speed, etc.
	// TODO: Offline caching weather response
	// TODO: Notification, if rain/thunderstorm coming.
	args := strings.Join(rq, ",")
	values := url.Values{}
	values.Add("latitude", strconv.FormatFloat(locationInfo.Latitude, 'f', -1, 64))
	values.Add("longitude", strconv.FormatFloat(locationInfo.Longitude, 'f', -1, 64))
	values.Add("current", args)
	fullUrl := network.WeatherEndpointApi + "forecast?" + values.Encode()

	resp, err := http.Get(fullUrl)
	if err != nil {
		return weatherResponse, locationInfo, err
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return weatherResponse, locationInfo, err
	}

	err = json.Unmarshal(bodyResp, &weatherResponse)
	if err != nil {
		return weatherResponse, locationInfo, err
	}
	return weatherResponse, locationInfo, nil
}

func GetHourlyWeather(location geocoding.LocationInfo) (HourlyWeatherResponse, error) {
	h := []string{"temperature_2m"}
	values := url.Values{}
	values.Add("latitude", strconv.FormatFloat(location.Latitude, 'f', -1, 64))
	values.Add("longitude", strconv.FormatFloat(location.Longitude, 'f', -1, 64))
	args := strings.Join(h, ",")
	values.Add("hourly", args)
	values.Add("past_days", "1")
	values.Add("forecast_days", "2")

	fullUrl := network.WeatherEndpointApi + "forecast?" + values.Encode()
	resp, err := http.Get(fullUrl)
	wtr := HourlyWeatherResponse{}
	if err != nil {
		return wtr, err
	}
	defer resp.Body.Close()

	bodyResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return wtr, err
	}

	err = json.Unmarshal(bodyResponse, &wtr)
	if err != nil {
		return wtr, err
	}
	return wtr, nil
}

func GetSunriseAndSunset(location geocoding.LocationInfo) (SunriseAndSunsetResponse, error) {
	/*
		Request generated according to API of sunrise and sunset of OpenMeteo
		Docs of Daily weather (including sunset & sunrise): https://open-meteo.com/en/docs#daily_weather_variables
	*/
	data := SunriseAndSunsetResponse{}
	values := url.Values{}
	values.Add("latitude", strconv.FormatFloat(location.Latitude, 'f', -1, 64))
	values.Add("longitude", strconv.FormatFloat(location.Longitude, 'f', -1, 64))
	values.Add("daily", "sunrise,sunset")
	values.Add("timezone", "auto")
	f := network.WeatherEndpointApi + "forecast?" + values.Encode()

	resp, err := http.Get(f)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	json.Unmarshal(b, &data)
	return data, nil
}

func GetCurrentWeatherName(weatherCode int) string {
	return weatherCodes[weatherCode]
}
