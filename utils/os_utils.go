package utils

import (
	"runtime"
	"os"
)

func GetRuntimeOs() string {
	return runtime.GOOS
}

func GetConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	return configDir
}