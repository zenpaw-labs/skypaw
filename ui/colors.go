package ui

import (
	"charm.land/lipgloss/v2"
	"github.com/zenpaw-labs/skypaw/utils/cfg"
)

var (
	// Ui colors
	ColorLocation         = lipgloss.Color("#40E0D0")
	ColorTemp             = lipgloss.Color("#FFD700")
	ColorWeatherName      = lipgloss.Color("#9bcbff")
	ColorTime             = lipgloss.Color("#5FFFD7")
	ColorDate             = lipgloss.Color("#98FB98")
	ColorHints            = lipgloss.Color("#4D4D4D")
	ColorError            = lipgloss.Color("#ff0000")
	ColorLoading          = lipgloss.Color("#ff7ad7")
	ColorOfflineOrTimeout = lipgloss.Color("#3f3f3f")

	// Weather colors
	ColorWeatherClear            = lipgloss.Color("#FFD700")
	ColorWeatherPartlyCloudy     = lipgloss.Color("#B0C4DE")
	ColorWeatherOvercast         = lipgloss.Color("#778899")
	ColorWeatherFog              = lipgloss.Color("#9DB4C0")
	ColorWeatherDrizzle          = lipgloss.Color("#87CEEB")
	ColorWeatherFreezingDrizzle  = lipgloss.Color("#B0E0E6")
	ColorWeatherRain             = lipgloss.Color("#4169E1")
	ColorWeatherFreezingRain     = lipgloss.Color("#00CED1")
	ColorWeatherSnowfall         = lipgloss.Color("#E0FFFF")
	ColorWeatherSnowGrains       = lipgloss.Color("#CAE9FF")
	ColorWeatherRainShowers      = lipgloss.Color("#1E90FF")
	ColorWeatherSnowShowers      = lipgloss.Color("#B0E0E6")
	ColorWeatherThunderstorm     = lipgloss.Color("#9400D3")
	ColorWeatherThunderstormHail = lipgloss.Color("#FF6347")
	ColorWeatherUnknown          = lipgloss.Color("#DA70D6")
)

func ApplyColors(c cfg.ColorsConfig) {
	ColorLocation = lipgloss.Color(c.Location)
	ColorTemp = lipgloss.Color(c.Temp)
	ColorWeatherName = lipgloss.Color(c.WeatherName)
	ColorTime = lipgloss.Color(c.Time)
	ColorDate = lipgloss.Color(c.Date)
	ColorHints = lipgloss.Color(c.Hints)
	ColorError = lipgloss.Color(c.Error)
	ColorLoading = lipgloss.Color(c.Loading)
	ColorOfflineOrTimeout = lipgloss.Color(c.OfflineOrTimeout)

	ColorWeatherClear = lipgloss.Color(c.WeatherClear)
	ColorWeatherPartlyCloudy = lipgloss.Color(c.WeatherPartlyCloudy)
	ColorWeatherOvercast = lipgloss.Color(c.WeatherOvercast)
	ColorWeatherFog = lipgloss.Color(c.WeatherFog)
	ColorWeatherDrizzle = lipgloss.Color(c.WeatherDrizzle)
	ColorWeatherFreezingDrizzle = lipgloss.Color(c.WeatherFreezingDrizzle)
	ColorWeatherRain = lipgloss.Color(c.WeatherRain)
	ColorWeatherFreezingRain = lipgloss.Color(c.WeatherFreezingRain)
	ColorWeatherSnowfall = lipgloss.Color(c.WeatherSnowfall)
	ColorWeatherSnowGrains = lipgloss.Color(c.WeatherSnowGrains)
	ColorWeatherRainShowers = lipgloss.Color(c.WeatherRainShowers)
	ColorWeatherSnowShowers = lipgloss.Color(c.WeatherSnowShowers)
	ColorWeatherThunderstorm = lipgloss.Color(c.WeatherThunderstorm)
	ColorWeatherThunderstormHail = lipgloss.Color(c.WeatherThunderstormHail)
	ColorWeatherUnknown = lipgloss.Color(c.WeatherUnknown)
}
