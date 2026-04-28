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
	return filepath.Join(configDir, "skypaw")
}

func GetBinaryDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return ex
}