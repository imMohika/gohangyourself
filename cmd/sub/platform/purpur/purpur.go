package purpur

import (
	"fmt"
	"github.com/imMohika/gohangyourself/api/purpur"
)

type PlatformPurpur struct {
}

func (p PlatformPurpur) Versions() ([]string, error) {
	versions, err := purpur.GetVersionList()
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func (p PlatformPurpur) LatestBuild(version string) (int, error) {
	latest, err := purpur.LatestBuild(version)
	return latest, err
}

func (p PlatformPurpur) DownloadURL(version string, build int) string {
	return fmt.Sprintf(
		"https://api.purpurmc.org/v2/purpur/%s/%d/download",
		version,
		build,
	)
}

func (p PlatformPurpur) FileName(version string, build int) string {
	return fmt.Sprintf("purpur-%s-%d.jar", version, build)
}
