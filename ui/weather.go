package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// TODO: Animation of weather
const (
	Clear = `
      \ | /
    -- ( ) --
      / | \
`

	PartlyCloudy = `
      \ | /
     --( ).--.
    .-(  )(   ).
   (___(___(___)
`

	Overcast = `
       .--.
    .-(    ).
   (___(___(__)
   (___________)
`

	Fog = `
    _ - _ - _ - _
     - _ - _ - _
    _ - _ - _ - _
`

	Drizzle = `
       .--.
    .-(    ).
   (___(___(__)
    ' ' ' '
`

	FreezingDrizzle = `
       .--.
    .-(    ).
   (___(___(__)
    * * * *
`

	Rain = `
       .--.
    .-(    ).
   (___(___(__)
    / / / /
   / / / /
`

	FreezingRain = `
       .--.
    .-(    ).
   (___(___(__)
    / * / *
   * / * /
`

	Snowfall = `
       .--.
    .-(    ).
   (___(___(__)
    * * * *
   * * * *
`

	SnowGrains = `
       .--.
    .-(    ).
   (___(___(__)
    . . . .
`

	RainShowers = `
       .--.
    .-(    ).
   (___(___(__)
   /// /// ///
`

	SnowShowers = `
       .--.
    .-(    ).
   (___(___(__)
   *** *** ***
`

	Thunderstorm = `
       .--.
    .-(    ).
   (___(___(__)
       /
      /
     /
`

	ThunderstormHail = `
       .--.
    .-(    ).
   (___(___(__)
     * / *
      / *
     * /
`

	Unknown = `
 _   _       _
| | | |_ __ | | ___ __   _____      ___ __
| | | | '_ \| |/ / '_ \ / _ \ \ /\ / / '_ \
| |_| | | | |   <| | | | (_) \ V  V /| | | |
 \___/|_| |_|_|\_\_| |_|\___/ \_/\_/ |_| |_|
`
)

type WeatherArt struct {
	Art   string
	Color lipgloss.Color
}

func GetCurrentWeatherArt(weatherCode int) WeatherArt {
	switch weatherCode {

	case 0, 1:
		return WeatherArt{
			Art:   Clear,
			Color: ColorWeatherClear,
		}
	case 2:
		return WeatherArt{
			Art:   PartlyCloudy,
			Color: ColorWeatherPartlyCloudy,
		}
	case 3:
		return WeatherArt{
			Art:   Overcast,
			Color: ColorWeatherOvercast,
		}
	case 45, 48:
		return WeatherArt{
			Art:   Fog,
			Color: ColorWeatherFog,
		}
	case 51, 53, 55:
		return WeatherArt{
			Art:   Drizzle,
			Color: ColorWeatherDrizzle,
		}
	case 56, 57:
		return WeatherArt{
			Art:   FreezingDrizzle,
			Color: ColorWeatherFreezingDrizzle,
		}
	case 61, 63, 65:
		return WeatherArt{
			Art:   Rain,
			Color: ColorWeatherRain,
		}
	case 66, 67:
		return WeatherArt{
			Art:   FreezingRain,
			Color: ColorWeatherFreezingRain,
		}
	case 71, 73, 75:
		return WeatherArt{
			Art:   Snowfall,
			Color: ColorWeatherSnowfall,
		}
	case 77:
		return WeatherArt{
			Art:   SnowGrains,
			Color: ColorWeatherSnowGrains,
		}
	case 80, 81, 82:
		return WeatherArt{
			Art:   RainShowers,
			Color: ColorWeatherRainShowers,
		}
	case 85, 86:
		return WeatherArt{
			Art:   SnowShowers,
			Color: ColorWeatherSnowShowers,
		}
	case 95:
		return WeatherArt{
			Art:   Thunderstorm,
			Color: ColorWeatherThunderstorm,
		}
	case 96, 99:
		return WeatherArt{
			Art:   ThunderstormHail,
			Color: ColorWeatherThunderstormHail,
		}
	default:
		return WeatherArt{
			Art:   Unknown,
			Color: ColorWeatherUnknown,
		}
	}
}
