package geocoding

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zenpaw-labs/skypaw/network"
	"github.com/zenpaw-labs/skypaw/utils/cfg"
)

//goland:noinspection GoNameStartsWithPackageName
type GeocodingResponse struct {
	Results        []LocationInfo `json:"results"`
	GenerationTime float64        `json:"generationtime_ms"`
}

type LocationInfo struct {
	City      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
	Region    string  `json:"admin1"`
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

type LocationProvider interface {
	URL() string
	Parse(body []byte) (LocationInfo, error)
}

type IpWhoProvider struct{}

func (p IpWhoProvider) Parse(body []byte) (LocationInfo, error) {
	var response struct {
		Country   string  `json:"country"`
		Region    string  `json:"region"`
		City      string  `json:"city"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return LocationInfo{}, err
	}
	return LocationInfo{
		Latitude:  response.Latitude,
		Longitude: response.Longitude,
		City:      response.City,
		Country:   response.Country,
		Region:    response.Region,
	}, nil
}

func (p IpWhoProvider) URL() string {
	return network.DetectLocationByIpWho
}

type IPAPIProvider struct{}

func (p IPAPIProvider) URL() string {
	return network.DetectLocationByNetworkIpApi
}

func (p IPAPIProvider) Parse(body []byte) (LocationInfo, error) {
	var response struct {
		/*
			Response from http://ip-api.com/json
		*/
		Country string  `json:"country"`
		Region  string  `json:"region"`
		City    string  `json:"city"`
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return LocationInfo{}, err
	}
	return LocationInfo{
		Latitude:  response.Lat,
		Longitude: response.Lon,
		City:      response.City,
		Country:   response.Country,
		Region:    response.Region,
	}, nil
}

type IPInfoProvider struct{}

func (p IPInfoProvider) Parse(body []byte) (LocationInfo, error) {
	var response struct {
		/*
			The struct data is under https://ipinfo.io/json response.
		*/
		City    string `json:"city"`
		Region  string `json:"region"`
		Country string `json:"country"`
		LOC     string `json:"loc"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return LocationInfo{}, err
	}

	l := strings.Split(response.LOC, ",")
	if len(l) < 2 {
		return LocationInfo{}, fmt.Errorf("invalid loc format: %s", response.LOC)
	}

	lat, _ := strconv.ParseFloat(l[0], 64)
	lon, _ := strconv.ParseFloat(l[1], 64)

	return LocationInfo{
		Latitude:  lat,
		Longitude: lon,
		City:      response.City,
		Country:   response.Country,
		Region:    response.Region,
	}, nil
}

func (p IPInfoProvider) URL() string {
	return network.DetectLocationByNetworkIpInfo
}

func FillLocationInfoFromCoords(l *LocationInfo) {
	// Used for PowerShell to fill geodata
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
		l.City = locData.Locality
	} else {
		l.City = locData.City
	}
	l.Country = locData.CountryName
	l.Region = locData.PrincipalSubdivision
}

func LocationDetectByNetwork(config cfg.UserConfig) (LocationInfo, error) {
	slog.Info("Detecting location by network with provider:", "provider_id", config.OptionalLocationProvider)
	var provider LocationProvider

	switch config.OptionalLocationProvider {
	case 1:
		provider = IpWhoProvider{}
	case 2:
		provider = IPAPIProvider{}
	case 3:
		provider = IPInfoProvider{}
	default:
		provider = IpWhoProvider{}
	}
	return executeDetection(provider)
}

func executeDetection(provider LocationProvider) (LocationInfo, error) {
	resp, err := http.Get(provider.URL())
	if err != nil {
		return LocationInfo{}, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationInfo{}, err
	}
	return provider.Parse(b)
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
