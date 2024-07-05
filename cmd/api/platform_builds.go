package api

import (
	"fmt"
	"github.com/imMohika/gohangyourself/api/hangar"
	"github.com/pterm/pterm"
)

func GetPlatformLatestBuild(platform, version string) (int, error) {
	spinner, _ := pterm.DefaultSpinner.Start("Getting builds")

	// todo: add spigot, etc
	builds, err := hangar.GetBuildList(platform, version)
	if err != nil {
		spinner.Fail()
		return 0, err
	}

	var latestBuild hangar.Builds
	for _, build := range builds {
		if build.Channel == "default" {
			latestBuild = build
			break
		}
	}

	if latestBuild.Build == 0 {
		return 0, fmt.Errorf("no stable build for version %s found :(", version)
	}

	spinner.Success()
	return latestBuild.Build, nil
}
