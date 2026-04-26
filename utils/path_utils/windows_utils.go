//go:build windows

package path_utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zenpaw-labs/skypaw/utils"
	"golang.org/x/sys/windows/registry"
)

func addToPath() error {
	err := addToWindowsPath(utils.GetBinaryDir())
	if err != nil {
		return err
	}
	return nil
}

func addToWindowsPath(dir string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry: %w", err)
	}
	defer k.Close()

	currentPath, _, err := k.GetStringValue("Path")
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}

	paths := strings.Split(currentPath, ";")
	for _, p := range paths {
		if strings.EqualFold(strings.TrimSpace(p), dir) {
			fmt.Println("Skypaw is already in the PATH!")
			return nil
		}
	}

	newPath := currentPath
	if newPath != "" && !strings.HasSuffix(newPath, ";") {
		newPath += ";"
	}
	newPath += dir

	err = k.SetExpandStringValue("Path", newPath)
	if err != nil {
		return fmt.Errorf("failed to update Path: %w", err)
	}

	fmt.Println("Successfully added to PATH!")
	return nil
}
