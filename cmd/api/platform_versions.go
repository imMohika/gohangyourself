package api

import (
	"github.com/imMohika/gohangyourself/api/hangar"
	"github.com/pterm/pterm"
)

func GetPlatformVersions(platform string) ([]string, error) {
	spinner, _ := pterm.DefaultSpinner.Start("Getting versions")
	defer spinner.Success()

	// todo: add spigot, etc
	versions, err := hangar.GetVersionList(platform)
	if err != nil {
		spinner.Fail()
		return nil, err
	}

	spinner.Success()
	return versions, nil
}
