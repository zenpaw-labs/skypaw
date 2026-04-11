package utils

import (
	"runtime"
)

func GetRuntimeOs() string {
	return runtime.GOOS
}
