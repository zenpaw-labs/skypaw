package location_utils

import (
	"github.com/zenpaw-labs/skypaw/network/geocoding"
)

func GetLocation(optionalProvider *int, winLocalDetect bool) (geocoding.LocationInfo, error) {
	coords, err := getLocationCoordinates(optionalProvider, winLocalDetect)
	if err != nil {
		return coords, err
	}
	return coords, nil
}
