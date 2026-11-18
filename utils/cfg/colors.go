package cfg

type ColorsConfig struct {
	Location         string `json:"location"`
	Temp             string `json:"temp"`
	WeatherName      string `json:"weather_name"`
	Time             string `json:"time"`
	Date             string `json:"date"`
	Hints            string `json:"hints"`
	Error            string `json:"error"`
	Loading          string `json:"loading"`
	OfflineOrTimeout string `json:"offline_or_timeout"`

	WeatherClear            string `json:"weather_clear"`
	WeatherPartlyCloudy     string `json:"weather_partly_cloudy"`
	WeatherOvercast         string `json:"weather_overcast"`
	WeatherFog              string `json:"weather_fog"`
	WeatherDrizzle          string `json:"weather_drizzle"`
	WeatherFreezingDrizzle  string `json:"weather_freezing_drizzle"`
	WeatherRain             string `json:"weather_rain"`
	WeatherFreezingRain     string `json:"weather_freezing_rain"`
	WeatherSnowfall         string `json:"weather_snowfall"`
	WeatherSnowGrains       string `json:"weather_snow_grains"`
	WeatherRainShowers      string `json:"weather_rain_showers"`
	WeatherSnowShowers      string `json:"weather_snow_showers"`
	WeatherThunderstorm     string `json:"weather_thunderstorm"`
	WeatherThunderstormHail string `json:"weather_thunderstorm_hail"`
	WeatherUnknown          string `json:"weather_unknown"`
}

func DefaultColors() ColorsConfig {
	return ColorsConfig{
		Location:         "#40E0D0",
		Temp:             "#FFD700",
		WeatherName:      "#9bcbff",
		Time:             "#5FFFD7",
		Date:             "#98FB98",
		Hints:            "#4D4D4D",
		Error:            "#ff0000",
		Loading:          "#ff7ad7",
		OfflineOrTimeout: "#3f3f3f",

		WeatherClear:            "#FFD700",
		WeatherPartlyCloudy:     "#B0C4DE",
		WeatherOvercast:         "#778899",
		WeatherFog:              "#9DB4C0",
		WeatherDrizzle:          "#87CEEB",
		WeatherFreezingDrizzle:  "#B0E0E6",
		WeatherRain:             "#4169E1",
		WeatherFreezingRain:     "#00CED1",
		WeatherSnowfall:         "#E0FFFF",
		WeatherSnowGrains:       "#CAE9FF",
		WeatherRainShowers:      "#1E90FF",
		WeatherSnowShowers:      "#B0E0E6",
		WeatherThunderstorm:     "#9400D3",
		WeatherThunderstormHail: "#FF6347",
		WeatherUnknown:          "#DA70D6",
	}
}
