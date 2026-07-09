package network

const (
	// Weather & Geocoding
	GeocodingEndpointApi string = "https://geocoding-api.open-meteo.com/v1/"
	WeatherEndpointApi   string = "https://api.open-meteo.com/v1/"
	// Location detectors
	DetectLocationByNetworkIpInfo string = "https://ipinfo.io/json"
	DetectLocationByNetworkIpApi  string = "http://ip-api.com/json"
	DetectLocationByIpWho         string = "https://ipwho.is"
	// For PowerShell only
	ReverseGeocodingApi string = "https://api.bigdatacloud.net/data/"
	// Links & autoupdater
	GithubLatestReleasePage     = "https://github.com/zenpaw-labs/skypaw/releases/latest"
	GithubLatestReleaseEndpoint = "https://api.github.com/repos/zenpaw-labs/skypaw/releases/latest"
)
