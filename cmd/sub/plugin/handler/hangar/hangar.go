package hangar

import (
	"github.com/imMohika/gohangyourself/api/hangar"
	"github.com/imMohika/gohangyourself/cmd/sub/plugin/internal"
	"strings"
)

type PluginHandler struct {
}

//func getPluginIDFromURL(url string) (string, error) {
//	url = strings.TrimSuffix(url, "/")
//	strings.Split(url, "/")
//}

func (h PluginHandler) FromURL(url string) (internal.PluginInfo, error) {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimSuffix(url, "/")
	parts := strings.Split(url, "/")
	name := len(parts) - 1

	return internal.PluginInfo{
		URL:  url,
		Name: parts[name],
	}, nil
}

// https://hangar.papermc.io/api/v1/projects/SayanVanish/versions
func (h PluginHandler) GetVersions(info internal.PluginInfo) ([]hangar.PluginVersion, error) {
	return hangar.GetPluginVersionList(info.Name)
}
