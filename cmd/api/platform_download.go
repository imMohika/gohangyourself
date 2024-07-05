package api

import (
	"github.com/imMohika/gohangyourself/api/hangar"
)

func DownloadPlatform(platform, version string, build int) error {
	// todo: add spigot, etc
	return hangar.DownloadPlatform(platform, version, build)
}
