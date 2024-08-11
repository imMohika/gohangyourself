package hangar

import (
	"fmt"
	"github.com/imMohika/gohangyourself/net"
	"slices"
	"time"
)

type Versions struct {
	Versions []string
}

func GetVersionList(platform string) ([]string, error) {
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s", platform)
	var versions Versions
	net.Get(url, "Couldn't get version list form hangar, platform="+platform, &versions)

	slices.Reverse(versions.Versions)
	return versions.Versions, nil
}

type BuildList struct {
	Builds []struct {
		Build    int       `json:"build"`
		Time     time.Time `json:"time"`
		Channel  string    `json:"channel"`
		Promoted bool      `json:"promoted"`
	}
}

func GetLatestBuild(platform string, version string) (int, error) {
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s/builds", platform, version)
	var buildList BuildList
	net.Get(url, "Failed to get latest build form hangar, platform="+platform, &buildList)

	builds := buildList.Builds
	slices.Reverse(builds)

	var latestBuild int
	for _, build := range builds {
		if build.Channel == "default" {
			latestBuild = build.Build
			break
		}
	}
	return latestBuild, nil
}
