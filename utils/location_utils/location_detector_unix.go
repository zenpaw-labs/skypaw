//go:build !windows

package location_utils

import (
	"skypaw/network/geocoding"
)

func getLocationCoordinates(osName string) (geocoding.LocationInfo, error) {
	return geocoding.LocationDetectByNetwork()
}
