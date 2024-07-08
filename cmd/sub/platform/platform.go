package platform

import (
	"errors"
	"flag"
	"github.com/imMohika/gohangyourself/cmd/download"
	"github.com/imMohika/gohangyourself/cmd/sub/platform/paper"
	"github.com/imMohika/gohangyourself/cmd/sub/platform/purpur"
	"github.com/imMohika/gohangyourself/cmd/sub/platform/velocity"
	"github.com/imMohika/gohangyourself/log"
	"github.com/pterm/pterm"
	"golang.org/x/exp/maps"
	"log/slog"
	"os"
	"slices"
	"strings"
)

type SubCommand struct{}

type Platform interface {
	Versions() ([]string, error)
	LatestBuild(version string) (int, error)
	DownloadURL(version string, build int) string
	FileName(version string, build int) string
}

var SupportedPlatforms = map[string]Platform{
	"paper":    paper.PlatformPaper{},
	"purpur":   purpur.PlatformPurpur{},
	"velocity": velocity.PlatformVelocity{},
}

var SupportedPlatformsKeys = maps.Keys(SupportedPlatforms)

func (p SubCommand) Handle(args []string) {
	flags := flag.NewFlagSet("platform", flag.ExitOnError)
	var platformFlag string
	flags.StringVar(&platformFlag, "platform", "", strings.Join(maps.Keys(SupportedPlatforms), ", "))
	flags.StringVar(&platformFlag, "p", "", strings.Join(maps.Keys(SupportedPlatforms), ", "))

	var versionFlag string
	flags.StringVar(&versionFlag, "version", "", "latest, 1.20.4")
	flags.StringVar(&versionFlag, "v", "", "latest, 1.20.4")

	err := flags.Parse(args)
	if err != nil {
		slog.Error("error parsing flags", "err", err)
	}

	platform, err := getPlatform(platformFlag)
	log.Error(err, "Couldn't get platform", platformFlag)
	if err != nil {
		slog.Error("Couldn't get platform",
			"args", args,
			"err", err)
		os.Exit(1)
	}

	version, err := getVersion(versionFlag, platform)
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

	url := platform.DownloadURL(version, build)
	fileName := platform.FileName(version, build)
	err = download.FromURL(url, fileName)
	if err != nil {
		slog.Error("Download failed",
			"platform", platform,
			"version", version,
			"build", build,
			"err", err)
		os.Exit(1)
	}
}

func getPlatform(platformFlag string) (Platform, error) {
	if platformFlag != "" {
		platform := strings.ToLower(platformFlag)
		if !slices.Contains(SupportedPlatformsKeys, platform) {
			return nil, errors.New("platform not supported")
		}

		return SupportedPlatforms[platform], nil
	}

	platform, err := pterm.DefaultInteractiveSelect.WithOptions(SupportedPlatformsKeys).Show("Select platform:")
	if err != nil {
		return nil, err
	}
	return SupportedPlatforms[platform], nil
}

func getVersion(versionFlag string, platform Platform) (string, error) {
	spinner, _ := pterm.DefaultSpinner.Start("Getting versions")
	versions, err := platform.Versions()
	if err != nil {
		spinner.Fail()
		return "", err
	}
	spinner.Success()

	if versionFlag != "" {
		version := strings.ToLower(versionFlag)
		if version == "latest" {
			return versions[0], nil
		}

		if !slices.Contains(versions, version) {
			return "", errors.New("invalid version")
		}
		return version, nil
	}

	version, err := pterm.DefaultInteractiveSelect.WithOptions(versions).Show("Select version")
	if err != nil {
		return "", err
	}
	return version, nil
}

func getBuild(platform Platform, version string) (int, error) {
	spinner, _ := pterm.DefaultSpinner.Start("Getting builds")
	latest, err := platform.LatestBuild(version)
	if err != nil {
		spinner.Fail()
		return 0, err
	}

	spinner.Success()
	return latest, nil
}
