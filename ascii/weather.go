package ascii

const (
	Clear = `
      \ | /
     - ( ) -
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

	Snow = `
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
)

func GetCurrentWeatherArt(weather string) string {
	return ""
}
