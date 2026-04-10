package main

import (
	"fmt"
	"skypaw/ascii"
	"skypaw/network/weather"
)

func main() {
	wr, _ := weather.GetCurrentWeather("Okhtyrka")
	fmt.Println(ascii.GetCurrentWeatherArt(wr.CurrentWeather.WeatherCode))
	fmt.Println(wr.CurrentWeather.Temperature2m, "°C")
	fmt.Println(weather.GetCurrentWeatherName(wr.CurrentWeather.WeatherCode))
	fmt.Scanln()
}
