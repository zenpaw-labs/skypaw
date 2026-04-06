package main

import (
	"fmt"
	"skypaw/network/geocoding"
	"skypaw/network/weather"
)

func main() {
	location := geocoding.SearchLocation("Okhtyrka")
	fmt.Println(weather.SearchWeather(location))
}
