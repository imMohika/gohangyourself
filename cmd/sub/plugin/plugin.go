package plugin

import (
	"fmt"
	api "github.com/imMohika/gohangyourself/api/hangar"
	"github.com/imMohika/gohangyourself/cmd/download"
	"github.com/imMohika/gohangyourself/cmd/sub/plugin/handler/hangar"
	"github.com/imMohika/gohangyourself/log"
	"github.com/pterm/pterm"
	"strconv"
	"strings"
)

type SubCommand struct {
}

func (s SubCommand) Handle(args []string) {
	url := args[0]
	pluginHandler := hangar.PluginHandler{}
	// handler, err := internal.GetHandlerFromURL(url)
	//log.Error(err, "Couldn't get plugin info", "url", url)
	pluginInfo, err := pluginHandler.FromURL(url)
	log.Error(err, "no meow")

	versions, err := pluginHandler.GetVersions(pluginInfo)
	if err != nil {
		return
	}

	version := selectVersion(versions)
	platform := selectPlatform(version)

	err = download.FromURL(platform.DownloadURL, platform.FileName)
	log.Error(err, "download failed")
}

func selectVersion(versions []api.PluginVersion) api.PluginVersion {
	options := make([]string, 0)
	for i, version := range versions {
		log.Info("here", "v", version, "i", i)
		options = append(options, fmt.Sprintf("%d - %s (%s)", i+1, version.Name, version.Channel.Name))
	}
	version, err := pterm.DefaultInteractiveSelect.WithOptions(options).Show("Select Version: ")
	log.Error(err, "idk")
	parts := strings.SplitN(version, " - ", 2)
	index, err := strconv.Atoi(parts[0])
	selectedVersion := versions[index-1]
	return selectedVersion
}

func selectPlatform(version api.PluginVersion) api.PluginPlatformDownload {
	options := make([]string, 0)
	for key, platform := range version.Downloads {
		options = append(options, fmt.Sprintf("%s (%s)", platform.Platform, version.PlatformDeps[key]))
	}
	platform, err := pterm.DefaultInteractiveSelect.WithOptions(options).Show("Select Platform: ")
	log.Error(err, "idk")
	parts := strings.SplitN(platform, " ", 2)
	selectedPlatform := version.Downloads[parts[0]]
	return selectedPlatform
}
