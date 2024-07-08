package velocity

import (
	"fmt"
	"github.com/imMohika/gohangyourself/api/hangar"
)

type PlatformVelocity struct{}

func (p PlatformVelocity) Versions() ([]string, error) {
	versions, err := hangar.GetVersionList("velocity")
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func (p PlatformVelocity) LatestBuild(version string) (int, error) {
	builds, err := hangar.GetBuildList("velocity", version)
	if err != nil {
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

	return latestBuild.Build, nil
}

func (p PlatformVelocity) DownloadURL(version string, build int) string {
	jarName := fmt.Sprintf("velocity-%s-%d.jar", version, build)
	return fmt.Sprintf("https://api.papermc.io/v2/projects/velocity/versions/%s/builds/%d/downloads/%s", version, build, jarName)
}

func (p PlatformVelocity) FileName(version string, build int) string {
	return fmt.Sprintf("velocity-%s-%d.jar", version, build)
}
