package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/masterminds/semver"
	"github.com/zenpaw-labs/skypaw/network"
)

func UnmarshalGithubLatestReleaseResponse(data []byte) (GithubLatestReleaseResponse, error) {
	var r GithubLatestReleaseResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GithubLatestReleaseResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GithubLatestReleaseResponse struct {
	TagName string  `json:"tag_name"`
	HTMLURL string  `json:"html_url"`
	Body    string  `json:"body"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func IsUpdatesAvailable(currentVersion string) (bool, string, error) {
	cVer, err := semver.NewVersion(currentVersion)
	if err != nil {
		return false, "", err
	}
	httpsResp, err := http.Get(network.GithubLatestReleaseEndpoint)
	if err != nil {
		return false, "", err
	}
	defer httpsResp.Body.Close()
	githubResponse := GithubLatestReleaseResponse{}
	b, err := io.ReadAll(httpsResp.Body)
	if err != nil {
		return false, "", err
	}
	json.Unmarshal(b, &githubResponse)
	lVer, err := semver.NewVersion(githubResponse.TagName)
	if err != nil {
		return false, "", err
	}

	return lVer.GreaterThan(cVer), lVer.String(), nil
}
