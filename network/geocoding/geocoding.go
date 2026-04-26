package geocoding

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zenpaw-labs/skypaw/network"
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

type IPAPIResponse struct {
	/*
		Response from http://ip-api.com/json
	*/
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

type IPInfoResponse struct {
	/*
		The struct data is under https://ipinfo.io/json response.
	*/
	IP       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	LOC      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
	Readme   string `json:"readme"`
}

type BigDataResponse struct {
	/*
		Response from https://api.bigdatacloud.net/data/
	*/
	City                 string `json:"city"`
	Locality             string `json:"locality"`
	CountryName          string `json:"countryName"`
	PrincipalSubdivision string `json:"principalSubdivision"`
}

func SearchLocation(name string) LocationInfo {
	/*
		Request generated according to Geocoding API of OpenMeteo.
		Docs of Geocoding API: https://open-meteo.com/en/docs/geocoding-api
	*/
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

func FillLocationInfoFromCoords(l *LocationInfo) {
	v := url.Values{}
	v.Add("latitude", strconv.FormatFloat(l.Latitude, 'f', -1, 64))
	v.Add("longitude", strconv.FormatFloat(l.Longitude, 'f', -1, 64))
	fullUrl := network.ReverseGeocodingApi + "reverse-geocode-client?" + v.Encode()
	resp, err := http.Get(fullUrl)
	if err != nil {
		return
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	locData := BigDataResponse{}
	err = json.Unmarshal(b, &locData)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if locData.Locality != "" {
		l.Name = locData.Locality
	} else {
		l.Name = locData.City
	}
	l.Country = locData.CountryName
	l.Admin1 = locData.PrincipalSubdivision
}

func locationDetectByNetworkIpApi() (LocationInfo, error) {
	var (
		locationInfo = LocationInfo{}
		response     = IPAPIResponse{}
	)
	resp, err := http.Get(network.DetectLocationByNetworkIpApi)
	if err != nil {
		return locationInfo, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return locationInfo, err
	}
	err = json.Unmarshal(b, &response)
	if err != nil {
		return locationInfo, err
	}
	locationInfo.Latitude = response.Lat
	locationInfo.Longitude = response.Lon
	locationInfo.Name = response.City
	locationInfo.Country = response.Country
	locationInfo.Admin1 = response.Region
	return locationInfo, nil
}

func locationDetectByNetworkIpInfo() (LocationInfo, error) {
	var (
		locationInfo = LocationInfo{}
		response     = IPInfoResponse{}
	)

	resp, err := http.Get(network.DetectLocationByNetworkIpInfo)
	if err != nil {
		return locationInfo, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return locationInfo, err
	}
	err = json.Unmarshal(b, &response)
	if err != nil {
		return locationInfo, err
	}
	l := strings.Split(response.LOC, ",")
	lat := l[0]
	lon := l[1]
	locationInfo.Latitude, _ = strconv.ParseFloat(lat, 64)
	locationInfo.Longitude, _ = strconv.ParseFloat(lon, 64)
	locationInfo.Name = response.City
	locationInfo.Country = response.Country
	locationInfo.Admin1 = response.Region
	return locationInfo, nil
}
func LocationDetectByNetwork(optionalProvider *int) (LocationInfo, error) {
	switch *optionalProvider {
	case 1:
		return locationDetectByNetworkIpApi()
	case 2:
		return locationDetectByNetworkIpInfo()
	default:
		return locationDetectByNetworkIpApi()
	}
}
