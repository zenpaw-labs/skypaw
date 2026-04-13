//go:build windows

package path_utils

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func addToPath(sys string) error {
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

	divider := ""
	if !strings.HasSuffix(oldPath, ";") {
		divider = ";"
	}

	newPath := oldPath + divider + targetDir + ";"
	err = k.SetStringValue("Path", newPath)
	if err != nil {
		return err
	}
	return nil
}
