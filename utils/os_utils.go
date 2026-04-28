package utils

import (
	"os"
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
	return configDir
}

func GetBinaryDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return ex
}