package paper

import (
	"fmt"
	"github.com/imMohika/gohangyourself/api/hangar"
)

type PlatformPaper struct {
}

func (p PlatformPaper) Versions() ([]string, error) {
	versions, err := hangar.GetVersionList("paper")
	if err != nil {
		return nil, err
	}
	return versions, nil
}
func (p PlatformPaper) LatestBuild(version string) (int, error) {
	build, err := hangar.GetLatestBuild("paper", version)
	if err != nil {
		return 0, err
	}

	if build == 0 {
		return 0, fmt.Errorf("no stable build for version %s found :(", version)
	}

	return build, nil
}

func (p PlatformPaper) FileName(version string, build int) string {
	return fmt.Sprintf("paper-%s-%d.jar", version, build)
}

func (p PlatformPaper) DownloadURL(version string, build int) string {
	return fmt.Sprintf("https://api.papermc.io/v2/projects/paper/versions/%s/builds/%d/downloads/%s", version, build, p.FileName(version, build))
}
