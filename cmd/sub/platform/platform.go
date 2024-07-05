package platform

import (
	"errors"
	"github.com/imMohika/gohangyourself/cmd/api"
	"github.com/pterm/pterm"
	"log/slog"
	"os"
	"slices"
	"strings"
)

type SubCommand struct{}

func (p SubCommand) Handle(args []string) {
	platform, err := getPlatform(args)
	if err != nil {
		slog.Error("Couldn't get platform",
			"args", args,
			"err", err)
		os.Exit(1)
	}

	version, err := getVersion(args, platform)
	if err != nil {
		slog.Error("Couldn't get version",
			"platform", platform,
			"args", args,
			"err", err)
		os.Exit(1)
	}

	build, err := getBuild(platform, version)
	if err != nil {
		slog.Error("Couldn't get latest build",
			"platform", platform,
			"version", version,
			"args", args,
			"err", err)
		os.Exit(1)
	}

	err = api.DownloadPlatform(platform, version, build)
	if err != nil {
		slog.Error("Download failed",
			"platform", platform,
			"version", version,
			"build", build,
			"err", err)
		os.Exit(1)
	}
}

var SupportedPlatforms = []string{"paper", "velocity"}

func getPlatform(args []string) (string, error) {
	if len(args) < 1 {
		platform, err := pterm.DefaultInteractiveSelect.WithOptions(SupportedPlatforms).Show("Select platform")
		if err != nil {
			return "", err
		}
		return platform, nil
	}

	platform := strings.ToLower(args[0])
	if !slices.Contains(SupportedPlatforms, platform) {
		return "", errors.New("platform not supported")
	}
	return platform, nil
}

func getVersion(args []string, platform string) (string, error) {
	versions, err := api.GetPlatformVersions(platform)
	if err != nil {
		return "", err
	}

	if len(args) < 2 {
		version, err := pterm.DefaultInteractiveSelect.WithOptions(versions).Show("Select version")
		if err != nil {
			return "", err
		}
		return version, nil
	}

	version := strings.ToLower(args[1])
	if version == "latest" {
		return versions[0], nil
	}

	if !slices.Contains(versions, version) {
		return "", errors.New("invalid version")
	}
	return version, nil
}

func getBuild(platform string, version string) (int, error) {
	latest, err := api.GetPlatformLatestBuild(platform, version)
	if err != nil {
		return 0, err
	}

	return latest, nil
}
