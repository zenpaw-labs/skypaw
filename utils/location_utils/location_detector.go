package location_utils

import (
	"github.com/zenpaw-labs/skypaw/network/geocoding"
)

func GetLocation(optionalProvider *int) (geocoding.LocationInfo, error) {
	coords, err := getLocationCoordinates(optionalProvider)
	if err != nil {
		return coords, err
	}
	return coords, nil
}
