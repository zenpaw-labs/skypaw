//go:build !windows

package location_utils

import (
	"github.com/zenpaw-labs/skypaw/network/geocoding"
)

func getLocationCoordinates() (geocoding.LocationInfo, error) {
	return geocoding.LocationDetectByNetwork()
}
