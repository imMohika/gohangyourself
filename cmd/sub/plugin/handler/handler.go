package handler

import (
	"fmt"
	"github.com/imMohika/gohangyourself/cmd/sub/plugin/handler/hangar"
	"github.com/imMohika/gohangyourself/cmd/sub/plugin/handler/modrinth"
	"github.com/imMohika/gohangyourself/cmd/sub/plugin/internal"
	"strings"
)

type PluginHandler interface {
	Name() string
	String() string
	LatestVersion()
	GetMeta() (internal.PluginMeta, error)
	GetVersionList() ([]internal.PluginVersion, error)
}

func GetHandler(url string) (PluginHandler, error) {
	var source string
	if strings.HasPrefix(url, "https") {
		trimmed := strings.TrimPrefix(url, "https://")

		switch {
		case strings.HasPrefix(trimmed, "modrinth.com"):
			source = "modrinth"
		case strings.HasPrefix(trimmed, "hangar.papermc.io"):
			source = "hangar"
		}
	} else {
		parts := strings.Split(url, ":")
		source = strings.TrimSuffix(parts[0], ":")
	}

	if source == "" {
		return nil, fmt.Errorf("plugin url must start with `hangar:`, `modrinth:`, or `https://` (%s)", url)
	}

	switch source {
	case "modrinth":
		return modrinth.FromURL(url), nil
	case "hangar":
		return hangar.FromURL(url), nil
	default:
		return nil, fmt.Errorf("%q is not a supported source", source)
	}
}
