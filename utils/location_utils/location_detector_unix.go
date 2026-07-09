//go:build !windows

package location_utils

import (
	"github.com/zenpaw-labs/skypaw/network/geocoding"
	"github.com/zenpaw-labs/skypaw/utils/cfg"
)

func getLocationCoordinates(config cfg.UserConfig, winLocalDetect bool) (geocoding.LocationInfo, error) {
	return geocoding.LocationDetectByNetwork(config)
}
