//go:build !windows

package path_utils

func addToPath(sys string) error {
	switch sys {

	case "darwin":
		err := addToPathDarwin()
		if err != nil {
			return err
		}
		return nil

	case "linux":
		err := addToPathLinux()
		if err != nil {
			return err
		}
	}
	return nil
}

func addToPathDarwin() error {
	return nil
}

func addToPathLinux() error {
	return nil
}
