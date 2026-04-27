package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

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
	URL             string    `json:"url"`
	AssetsURL       string    `json:"assets_url"`
	UploadURL       string    `json:"upload_url"`
	HTMLURL         string    `json:"html_url"`
	ID              int64     `json:"id"`
	Author          Author    `json:"author"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Immutable       bool      `json:"immutable"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []Asset   `json:"assets"`
	TarballURL      string    `json:"tarball_url"`
	ZipballURL      string    `json:"zipball_url"`
	Body            string    `json:"body"`
}

type Asset struct {
	URL                string      `json:"url"`
	ID                 int64       `json:"id"`
	NodeID             string      `json:"node_id"`
	Name               string      `json:"name"`
	Label              interface{} `json:"label"`
	Uploader           Author      `json:"uploader"`
	ContentType        string      `json:"content_type"`
	State              string      `json:"state"`
	Size               int64       `json:"size"`
	Digest             string      `json:"digest"`
	DownloadCount      int64       `json:"download_count"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	BrowserDownloadURL string      `json:"browser_download_url"`
}

type Author struct {
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	UserViewType      string `json:"user_view_type"`
	SiteAdmin         bool   `json:"site_admin"`
}

func IsUpdatesAvailable(currentVersion string) (bool, string, error) {
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
	cVer, err := semver.NewVersion(currentVersion)
	if err != nil {
		return false, "", err
	}
	lVer, err := semver.NewVersion(githubResponse.TagName)
	if err != nil {
		return false, "", err
	}

	return lVer.GreaterThan(cVer), lVer.String(), nil
}
