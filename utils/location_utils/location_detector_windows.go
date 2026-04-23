//go:build windows

package location_utils

import (
	"bytes"
	"encoding/json"
	"os/exec"

	"github.com/zenpaw-labs/skypaw/network/geocoding"
)

func getLocationCoordinates() (geocoding.LocationInfo, error) {
	location, err := locationDetectorByPS()
	if err != nil {
		return geocoding.LocationDetectByNetwork()
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
		return geocoding.LocationDetectByNetwork()
	}
	return locationInfo, nil
}
