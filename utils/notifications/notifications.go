package notifications

import (
	"fmt"
	"slices"
	"sync"
	"time"

	"charm.land/log/v2"
	"github.com/gen2brain/beeep"
	"github.com/zenpaw-labs/skypaw/network/weather"
	"github.com/zenpaw-labs/skypaw/utils/cfg"
)

const (
	cooldown = 3 * time.Hour
)

var (
	// cooldown
	state   notifyState
	stateMu sync.Mutex

	badWeatherCodes = []int{
		/*
			Thunderstorm codes
		*/
		95, 96, 99,

		/*
			Rain codes
		*/
		61, 63, 65, 66, 67, 80, 81, 82,

		/*
			Snow codes
		*/
		71, 73, 75, 77, 85, 86,

		/*
			Freezing drizzle codes
		*/
		56, 57, 1, 2,
	}
)

type notifyState struct {
	lastHour time.Time
	lastCode int
	lastSent time.Time
}

func BadWeatherNotification(config cfg.UserConfig, wtr weather.HourlyWeatherResponse, city string) {
	next := time.Now().Truncate(time.Hour).Add(time.Hour)
	limit := next.Add(time.Duration(config.BadWeatherLookahead))

	for i, code := range wtr.Hourly.WeatherCode {
		t, err := time.ParseInLocation("2006-01-02T15:04", wtr.Hourly.Time[i], time.UTC)
		if err != nil || t.Before(next) {
			continue
		}
		if t.After(limit) {
			break
		}

		if slices.Contains(badWeatherCodes, code) {
			if shouldNotify(t, code) {
				log.Debug("Bad Weather Notification", "weather_code", code)
				err := beeep.Alert(
					"Bad Weather Notice",
					fmt.Sprintf("%s is expected in %s at %s", weather.GetCurrentWeatherName(code), city, t.In(time.Local).Format("15:04")),
					"",
				)
				if err != nil {
					log.Error("Failed to send 'Bad Weather Notice'", "error", err, "weather_code", code)
				} else {
					markNotified(t, code)
				}
			}
			break
		}
	}
}

func shouldNotify(t time.Time, code int) bool {
	stateMu.Lock()
	defer stateMu.Unlock()

	if t.Equal(state.lastHour) && code == state.lastCode {
		return false
	}
	/*
		There are conditions, that can dodge cooldown

		if slices.Contains(thunderCodes, code) && !slices.Contains(thunderCodes, state.lastCode) {
			return true
		}
	*/
	return time.Since(state.lastSent) >= cooldown
}

func markNotified(t time.Time, code int) {
	stateMu.Lock()
	defer stateMu.Unlock()
	state = notifyState{lastHour: t, lastCode: code, lastSent: time.Now()}
}
