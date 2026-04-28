package ascii

const (
	Clear = `  \ | /
 - ( ) -
  / | \`

	PartlyCloudy = `   \ | / 
  --( ).--.
 .-(  )(   ).
(___(___(___)`

	Overcast = `    .--.
 .-(    ).
(___(___(__)
(___________)`

	Fog = ` _ - _ - _ - _
  - _ - _ - _
 _ - _ - _ - _`

	Drizzle = `    .--.
 .-(    ).
(___(___(__)
 ' ' ' '`

	FreezingDrizzle = `    .--.
 .-(    ).
(___(___(__)
 * * * *`

	Rain = `    .--.
 .-(    ).
(___(___(__)
 / / / /
/ / / /`

	FreezingRain = `    .--.
 .-(    ).
(___(___(__)
 / * / *
* / * /`

	Snowfall = `    .--.
 .-(    ).
(___(___(__)
 * * * *
* * * *`

	SnowGrains = `    .--.
 .-(    ).
(___(___(__)
 . . . .`

	RainShowers = `    .--.
 .-(    ).
(___(___(__)
/// /// ///`

	SnowShowers = `    .--.
 .-(    ).
(___(___(__)
*** *** ***`

	Thunderstorm = `    .--.
 .-(    ).
(___(___(__)
    /
   /
  /`

	ThunderstormHail = `    .--.
 .-(    ).
(___(___(__)
  * / *
   / *
  * /`

	Unknown = ` _   _       _                              
| | | |_ __ | | ___ __   _____      ___ __  
| | | | '_ \| |/ / '_ \ / _ \ \ /\ / / '_ \ 
| |_| | | | |   <| | | | (_) \ V  V /| | | |
 \___/|_| |_|_|\_\_| |_|\___/ \_/\_/ |_| |_|`
)

func GetCurrentWeatherArt(weatherCode int) string {
	switch weatherCode {
	case 0, 1:
		return Clear
	case 2:
		return PartlyCloudy
	case 3:
		return Overcast
	case 45, 48:
		return Fog
	case 51, 53, 55:
		return Drizzle
	case 56, 57:
		return FreezingDrizzle
	case 61, 63, 65:
		return Rain
	case 66, 67:
		return FreezingRain
	case 71, 73, 75:
		return Snowfall
	case 77:
		return SnowGrains
	case 80, 81, 82:
		return RainShowers
	case 85, 86:
		return SnowShowers
	case 95:
		return Thunderstorm
	case 96, 99:
		return ThunderstormHail
	default:
		return Unknown
	}
}
