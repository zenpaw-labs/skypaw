package path_utils

import (
	"skypaw/utils"
)

func AddToPath() error {
	osName := utils.GetRuntimeOs()
	switch osName {
	case "windows":
		err := addToPath(osName)
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

func darwin() error {
	return nil
}

func linux() error {
	return nil
}
