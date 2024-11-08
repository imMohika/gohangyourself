package download

import (
	"errors"
	"flag"
	"fmt"
	"github.com/imMohika/gohangyourself/cmd/download"
	"github.com/imMohika/gohangyourself/cmd/sub/plugin/handler"
	"github.com/imMohika/gohangyourself/cmd/sub/plugin/internal"
	"github.com/imMohika/gohangyourself/log"
	"github.com/pterm/pterm"
	"slices"
	"strconv"
	"strings"
)

type SubCommand struct {
}

func (s SubCommand) Handle(args []string) {
	flags := flag.NewFlagSet("plugin:download", flag.ExitOnError)
	var FlagLatestVersion bool
	flags.BoolVar(&FlagLatestVersion, "latest", false, "")
	flags.BoolVar(&FlagLatestVersion, "l", false, "")
	var FlagPlatform string
	flags.StringVar(&FlagPlatform, "platform", "", "")
	flags.StringVar(&FlagPlatform, "p", "", "")

	err := flags.Parse(args)
	if err != nil {
		log.Error(err, "error parsing flags")
	}

	if len(flags.Args()) != 1 {
		log.Warn("wrong usage. please provide url")
		return
	}

	url := strings.TrimSpace(flags.Args()[0])
	pluginHandler, err := handler.GetHandler(url)
	if err != nil {
		log.Fatal(err, "can not get handler", "url", url)
		return
	}

	versions, err := getVersions(pluginHandler)
	if err != nil {
		log.Fatal(err, "can not plugin versions", "url", url, "handler", pluginHandler.Name())
		return
	}

	version, err := selectVersion(versions, FlagLatestVersion, FlagPlatform)
	if err != nil {
		log.Fatal(err, "can not select plugin version", "url", url, "handler", pluginHandler.Name())
		return
	}

	file, err := selectPlatform(version, FlagPlatform)
	if err != nil {
		log.Fatal(err, "can not select plugin version file (select loader)", "url", url, "handler", pluginHandler.Name())
		return
	}

	err = download.FromURL(file.URL, file.Name)
	if err != nil {
		log.Fatal(err, "Failed to download plugin",
			"url", file.URL,
			"version", version.ID,
			"handler", pluginHandler.Name())
	}
}

func getVersions(pluginHandler handler.PluginHandler) ([]internal.PluginVersion, error) {
	spinner, _ := pterm.DefaultSpinner.Start("Getting versions")

	versions, err := pluginHandler.GetVersionList()
	if err != nil {
		return nil, err
	}

	spinner.Success()
	return versions, nil
}

func selectVersion(versions []internal.PluginVersion, latest bool, platform string) (internal.PluginVersion, error) {
	if latest {
		for _, version := range versions {
			if slices.Contains(version.Loaders, platform) {
				return version, nil
			}
		}

		return internal.PluginVersion{}, errors.New("could not find latest version for specified platform")
	}

	options := make([]string, len(versions))
	for i, version := range versions {
		options[i] = fmt.Sprintf("%d. %s [%s]",
			i+1, version.Name, strings.Join(version.Loaders, ","))
	}

	version, err := pterm.DefaultInteractiveSelect.WithOptions(options).Show("Please select a version:")
	if err != nil {
		return internal.PluginVersion{}, err
	}

	parts := strings.Split(version, ".")
	idx, err := strconv.Atoi(parts[0])
	if err != nil {
		return internal.PluginVersion{}, err
	}

	return versions[idx-1], nil
}

func selectPlatform(version internal.PluginVersion, platform string) (internal.PluginFile, error) {
	if platform != "" {
		for _, file := range version.Files {
			if slices.Contains(file.Loaders, platform) {
				return file, nil
			}
		}
		return internal.PluginFile{}, errors.New("could not find specified platform")
	}

	if len(version.Files) == 1 {
		return version.Files[0], nil
	}

	options := make([]string, len(version.Files))
	for i, file := range version.Files {
		options[i] = fmt.Sprintf("%d. %s",
			i+1, strings.Join(file.Loaders, ", "))
	}

	file, err := pterm.DefaultInteractiveSelect.WithOptions(options).Show("Please select a loader:")
	if err != nil {
		return internal.PluginFile{}, err
	}

	parts := strings.Split(file, ".")
	idx, err := strconv.Atoi(parts[0])
	if err != nil {
		return internal.PluginFile{}, err
	}

	return version.Files[idx-1], nil
}
