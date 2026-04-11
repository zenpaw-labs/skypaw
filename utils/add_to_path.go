package utils

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func AddToPath() error {
	osName := GetRuntimeOs()
	switch osName {
	case "windows":
		err := windows()
		if err != nil {
			return err
		}
		return nil

	case "darwin":
		darwin()
		return nil

	case "linux":
		linux()
		return nil
	}
	return nil
}

func windows() error {
	path, err := os.Executable()
	if err != nil {
		return err
	}

	targetDir := filepath.Dir(path)
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	oldPath, _, err := k.GetStringValue("Path")
	if err != nil {
		return err
	}

	if strings.Contains(oldPath, targetDir) {
		return nil
	}

	newPath := oldPath + ";" + targetDir
	err = k.SetStringValue("Path", newPath)
	if err != nil {
		return err
	}
	return nil
}

func darwin() error {
	return nil
}

func linux() error {
	return nil
}
