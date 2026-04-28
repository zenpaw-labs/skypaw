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