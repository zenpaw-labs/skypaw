package main

import (
	"fmt"
	"skypaw/network/weather"
)

func main() {

	wr, _ := weather.GetCurrentWeather("Okhtyrka")
	fmt.Println("Weather: ", wr)
	fmt.Scanln()
}
