package location_utils

import (
	"github.com/zenpaw-labs/skypaw/network/geocoding"
)

func GetLocationFromOs() (geocoding.LocationInfo, error) {
	coords, err := getLocationCoordinates()
	if err != nil {
		return coords, err
	}
	locName, err := geocoding.GetLocationFromCoords(coords)
	if err != nil {
		return coords, err
	}
	coords.Name = locName.Name
	return coords, nil
}
