package purpur

import (
	"fmt"
	"github.com/imMohika/gohangyourself/log"
	"github.com/imMohika/gohangyourself/net"
	"slices"
	"strconv"
)

type Versions struct {
	Versions []string
}

func GetVersionList() ([]string, error) {
	url := "https://api.purpurmc.org/v2/purpur"
	var versions Versions
	net.Get(url, "Couldn't get version list for purpur", &versions)

	slices.Reverse(versions.Versions)
	return versions.Versions, nil
}

type BuildsList struct {
	Builds struct {
		All    []string
		Latest string
	}
}

func GetLatestBuild(version string) (int, error) {
	url := fmt.Sprintf("https://api.purpurmc.org/v2/purpur/%s", version)
	var buildsList BuildsList
	net.Get(url, "failed to get build list", &buildsList)
	latest, err := strconv.Atoi(buildsList.Builds.Latest)
	log.Error(err, "failed to get latest build")

	return latest, nil
}
