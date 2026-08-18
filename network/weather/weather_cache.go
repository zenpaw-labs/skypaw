package weather

import (
	"encoding/json"
	"os"
	"time"

	"github.com/zenpaw-labs/skypaw/network/geocoding"
	"github.com/zenpaw-labs/skypaw/utils"
	"github.com/zenpaw-labs/skypaw/utils/cfg"
)

type WeatherCache struct {
	LastUpdatedAt time.Time                `json:"last_updated_at"`
	Location      geocoding.LocationInfo   `json:"location"`
	Weather       WeatherResponse          `json:"weather"`
	Hourly        HourlyWeatherResponse    `json:"hourly"`
	SunriseSunset SunriseAndSunsetResponse `json:"sunrise_sunset"`
	UnitSystem    cfg.UnitSystem           `json:"unit_system"`
}

func Load(cfg cfg.UserConfig) (*WeatherCache, error) {
	path := utils.GetWeatherCacheDir()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cache WeatherCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	if cfg.Units != cache.UnitSystem {
		convertTemperatureUnits(&cache, cfg.Units)
	}
	return &cache, nil
}

func Save(l geocoding.LocationInfo, w WeatherResponse, h HourlyWeatherResponse, s SunriseAndSunsetResponse, u cfg.UnitSystem) error {
	path := utils.GetWeatherCacheDir()
	cache := WeatherCache{
		LastUpdatedAt: time.Now(),
		Location:      l,
		Weather:       w,
		Hourly:        h,
		SunriseSunset: s,
		UnitSystem:    u,
	}
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func convertTemperatureUnits(cache *WeatherCache, currentSystem cfg.UnitSystem) {
	if cache == nil || cache.UnitSystem == currentSystem {
		return
	}

	toImperial := currentSystem == cfg.Imperial

	convert := func(val float64) float64 {
		if toImperial {
			return val*9/5 + 32
		}
		return (val - 32) * 5 / 9
	}
	cache.Weather.CurrentWeather.Temperature2m = convert(cache.Weather.CurrentWeather.Temperature2m)

	for i, temp := range cache.Hourly.Hourly.Temperature2m {
		cache.Hourly.Hourly.Temperature2m[i] = convert(temp)
	}

	cache.UnitSystem = currentSystem
}
