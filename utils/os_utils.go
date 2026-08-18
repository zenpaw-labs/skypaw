package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetRuntimeOs() string {
	return runtime.GOOS
}

func GetConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	return filepath.Join(configDir, "Zenpaw Labs", "skypaw")
}

func GetConfigFile() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	path := filepath.Join(configDir, "Zenpaw Labs", "skypaw", "config.json")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", err
	}

	return path, nil
}

func GetWeatherCacheDir() string {
	cachedir, _ := os.UserConfigDir()
	path := filepath.Join(cachedir, "Zenpaw Labs", "skypaw", "last_weather_update.json")
	return path
}
