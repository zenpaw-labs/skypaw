//go:build windows

package location_utils

import (
	"bytes"
	"encoding/json"
	"os/exec"

	"charm.land/log/v2"

	"github.com/zenpaw-labs/skypaw/network/geocoding"
	"github.com/zenpaw-labs/skypaw/utils/cfg"
)

func getLocationCoordinates(config cfg.UserConfig, winLocalDetect bool) (geocoding.LocationInfo, error) {
	if !winLocalDetect {
		log.Info("Detecting location on Windows by optional provider (flag or config)", "provider_id", config.OptionalLocationProvider)
		return geocoding.LocationDetectByNetwork(config)
	}

	location, err := locationDetectorByPS()
	if err != nil || (location.Latitude == 0 && location.Longitude == 0) {
		log.Warn("Detecting location on Windows by optional provider (PS Error)", "provider_id", config.OptionalLocationProvider)
		return geocoding.LocationDetectByNetwork(config)
	}
	return location, err
}

func locationDetectorByPS() (geocoding.LocationInfo, error) {
	var (
		locationInfo geocoding.LocationInfo
	)

	psScript := `
	Add-Type -AssemblyName System.Device;
	$gw = New-Object System.Device.Location.GeoCoordinateWatcher;
	$gw.Start();
	while ($gw.Status -ne 'Ready' -and $gw.Permission -ne 'Denied') { Start-Sleep -Milliseconds 100 };
	$data = $gw.Position.Location | Select-Object Latitude, Longitude;
	$gw.Stop();
	$data | ConvertTo-Json
	`

	cmd := exec.Command("powershell", "-Command", psScript)

	var out bytes.Buffer
	cmd.Stdout = &out
	var errOut bytes.Buffer
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return locationInfo, err
	}

	err = json.Unmarshal(out.Bytes(), &locationInfo)
	if err != nil || (locationInfo.Latitude == 0 && locationInfo.Longitude == 0) {
		return locationInfo, err
	}
	geocoding.FillLocationInfoFromCoords(&locationInfo)
	return locationInfo, nil
}
