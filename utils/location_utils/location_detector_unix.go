//go:build !windows

package location_utils

import (
	"github.com/zenpaw-labs/skypaw/network/geocoding"
)

func getLocationCoordinates(osName string) (geocoding.LocationInfo, error) {
	return geocoding.LocationDetectByNetwork()
}
