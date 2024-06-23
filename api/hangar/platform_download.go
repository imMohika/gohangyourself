package hangar

import (
	"fmt"
	"gohangyourself/cmd/download"
)

func DownloadPlatform(platform, version string, build int) error {
	//     JAR_NAME=${PROJECT}-${MINECRAFT_VERSION}-${LATEST_BUILD}.jar
	jarName := fmt.Sprintf("%s-%s-%d.jar", platform, version, build)
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s/builds/%d/downloads/%s", platform, version, build, jarName)

	err := download.FromURL(url)
	if err != nil {
		return err
	}

	return nil
}
