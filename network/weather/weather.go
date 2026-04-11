package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"skypaw/network"
	"skypaw/network/geocoding"
	"strconv"
	"strings"
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

/*
	Requests generated for according to API Scheme of current weather by OpenMeteo
	Docs of current weather API: https://open-meteo.com/en/docs#current_weather
*/

func GetCurrentWeather(location string) (WeatherResponse, geocoding.LocationInfo, error) {
	locationInfo := geocoding.SearchLocation(location)
	weather, err := GetCurrentWeatherByLocationInfo(locationInfo)
	if err != nil {
		return weather, locationInfo, err
	}
	return weather, locationInfo, nil
}

func GetCurrentWeatherByLocationInfo(locationInfo geocoding.LocationInfo) (WeatherResponse, error) {

	var (
		weatherResponse WeatherResponse
	)
	rq := []string{"temperature_2m", "is_day", "weather_code", "wind_speed_10m"}
	args := strings.Join(rq, ",")
	values := url.Values{}
	values.Add("latitude", strconv.FormatFloat(locationInfo.Latitude, 'f', -1, 64))
	values.Add("longitude", strconv.FormatFloat(locationInfo.Longitude, 'f', -1, 64))
	values.Add("current", args)
	fullUrl := network.WeatherEndpointApi + "forecast?" + values.Encode()

	resp, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("An error occurred: ", err, ".")
		return weatherResponse, err
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("An error occurred: ", err, ".")
		return weatherResponse, err
	}

	err = json.Unmarshal(bodyResp, &weatherResponse)
	if err != nil {
		fmt.Println("An error occurred: ", err, ".")
	}
	return weatherResponse, nil
}

func GetCurrentWeatherFromCoordinates() string {
	// TODO
	return "sex"
}

func GetCurrentWeatherName(weatherCode int) string {
	return weatherCodes[weatherCode]
}
