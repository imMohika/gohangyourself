package hangar

import (
	"fmt"
	"github.com/imMohika/gohangyourself/net"
	"github.com/tidwall/gjson"
)

type PluginPlatformDownload struct {
	FileName    string
	DownloadURL string
	Platform    string
}

type PluginVersionChannel struct {
	Name string
}

type PluginVersion struct {
	Name         string
	Channel      PluginVersionChannel
	Downloads    map[string]PluginPlatformDownload
	PlatformDeps map[string]string
}

// https://hangar.papermc.io/api/v1/projects/SayanVanish/versions
func GetPluginVersionList(name string) ([]PluginVersion, error) {
	url := fmt.Sprintf("https://hangar.papermc.io/api/v1/projects/%s/versions", name)
	resp := net.GetGJSON(url, "Couldn't get version list form hangar, plugin="+name)
	versions := make([]PluginVersion, 0)
	resp.Get("result").ForEach(func(key, value gjson.Result) bool {
		downloads := make(map[string]PluginPlatformDownload)
		value.Get("downloads").ForEach(func(key, value gjson.Result) bool {
			downloads[key.String()] = PluginPlatformDownload{
				FileName:    value.Get("fileInfo.name").String(),
				DownloadURL: value.Get("downloadUrl").String(),
				Platform:    key.String(),
			}
			return true
		})

		platformDeps := make(map[string]string)
		value.Get("platformDependenciesFormatted").ForEach(func(key, value gjson.Result) bool {
			platformDeps[key.String()] = value.String()
			return true
		})

		versions = append(versions, PluginVersion{
			Name:         value.Get("name").String(),
			Channel:      PluginVersionChannel{Name: value.Get("channel.name").String()},
			Downloads:    downloads,
			PlatformDeps: platformDeps,
		})
		return true
	})

	return versions, nil
}
