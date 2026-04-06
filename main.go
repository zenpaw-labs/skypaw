package main

import (
	"fmt"
	"skypaw/network/geocoding"
	"skypaw/network/weather"
)

func main() {
	location := geocoding.SearchLocation("Okhtyrka")
	wr := weather.GetCurrentWeather(location)
	fmt.Println("Location: ", location)
	fmt.Println("Weather: ", wr)
	fmt.Scanln()
}
