package geocoding

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"skypaw/network"
)

//goland:noinspection GoNameStartsWithPackageName
type GeocodingResponse struct {
	Results        []LocationInfo `json:"results"`
	GenerationTime float64        `json:"generationtime_ms"`
}

type LocationInfo struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Elevation   float64 `json:"elevation"`
	FeatureCode string  `json:"feature_code"`
	CountryCode string  `json:"country_code"`
	Timezone    string  `json:"timezone"`
	Population  int     `json:"population"`
	CountryID   int     `json:"country_id"`
	Country     string  `json:"country"`
	Admin1      string  `json:"admin1"`
	Admin2      string  `json:"admin2"`
	Admin3      string  `json:"admin3"`
	Admin4      string  `json:"admin4"`
	Admin1ID    int     `json:"admin1_id"`
	Admin2ID    int     `json:"admin2_id"`
	Admin3ID    int     `json:"admin3_id"`
	Admin4ID    int     `json:"admin4_id"`
}

/*
	Request generated according to Geocoding API of OpenMeteo.
	Docs of Geocoding API: https://open-meteo.com/en/docs/geocoding-api
*/

func SearchLocation(name string) LocationInfo {
	var (
		locatonInfo LocationInfo
		geoData     GeocodingResponse
	)

	params := url.Values{}
	params.Add("name", name)
	fullUrl := network.GeocodingEndpointApi + "search?" + params.Encode()

	resp, err := http.Get(fullUrl)
	if err != nil {
		fmt.Println("An error occurred: ", err)
		return locatonInfo
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("An error occurred: ", err)
		return locatonInfo
	}

	err = json.Unmarshal(response, &geoData)
	if err != nil {
		fmt.Println("An error occurred: ", err)
		return locatonInfo
	}
	if len(geoData.Results) > 0 {
		locatonInfo = geoData.Results[0]
	}

	return locatonInfo
}
