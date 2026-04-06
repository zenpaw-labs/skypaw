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

func GetCurrentWeather(locationInfo geocoding.LocationInfo) WeatherResponse {
	/*
		Request generated for according to API Scheme of current weather by OpenMeteo
		Docs of current weather API: https://open-meteo.com/en/docs#current_weather
	*/
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
		return weatherResponse
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("An error occurred: ", err, ".")
		return weatherResponse
	}

	err = json.Unmarshal(bodyResp, &weatherResponse)
	if err != nil {
		fmt.Println("An error occurred: ", err, ".")
	}
	return weatherResponse
}
