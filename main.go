package main

import (
	"fmt"
	"skypaw/network/weather"
)

func main() {
	wr, _ := weather.GetCurrentWeather("Okhtyrka")
	fmt.Println("Weather: ", weather.GetCurrentWeatherName(wr.CurrentWeather.WeatherCode))
	fmt.Scanln()
}
