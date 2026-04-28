//go:build !windows

package path_utils

import (
	"fmt"
	"os"

)

func addToPath() error {
	return installUnix(getBinaryDir())
}

func installUnix(path string) error {
	targetPath := "/usr/local/bin/skypaw"

	input, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Could not read binary: %w", err)
	}

	err = os.WriteFile(targetPath, input, 0755)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("Permission denied. Please run: sudo ./skypaw --install")
		}
		return fmt.Errorf("Could not install: %w", err)
	}

	fmt.Printf("Successfully installed to %s.\n", targetPath)
	return nil
}

func getBinaryDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return ex
}