package location_utils

import "skypaw/network/geocoding"

func GetLocationFromOs(osName string) (geocoding.LocationInfo, error) {
	return getLocationCoordinates(osName)
}
