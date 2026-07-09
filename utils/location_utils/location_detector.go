package location_utils

import (
	"github.com/zenpaw-labs/skypaw/network/geocoding"
	"github.com/zenpaw-labs/skypaw/utils/cfg"
)

func GetLocation(config cfg.UserConfig, winLocalDetect bool) (geocoding.LocationInfo, error) {
	coords, err := getLocationCoordinates(config, winLocalDetect)
	if err != nil {
		return coords, err
	}
	return coords, nil
}
