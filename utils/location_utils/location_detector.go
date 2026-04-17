package location_utils

import "github.com/zenpaw-labs/skypaw/network/geocoding"

func GetLocationFromOs(osName string) (geocoding.LocationInfo, error) {
	return getLocationCoordinates(osName)
}
