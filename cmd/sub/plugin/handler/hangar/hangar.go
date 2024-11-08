package hangar

import (
	"fmt"
	"github.com/imMohika/gohangyourself/cmd/sub/plugin/internal"
	"github.com/imMohika/gohangyourself/log"
	"github.com/imMohika/gohangyourself/net"
	"regexp"
	"strings"
)

type PluginHandler struct {
	url  string
	slug string
}

var hangarRegex = regexp.MustCompile(`^https:\/\/hangar\.papermc\.io\/.+\/(.+)$`)

func FromURL(url string) *PluginHandler {
	if strings.HasPrefix(url, "hangar:") {
		_, slug, found := strings.Cut(url, ":")
		if !found {
			log.FatalMsg("invalid url passed to plugin.handler.hangar", "url", url)
		}
		return &PluginHandler{
			url:  url,
			slug: slug,
		}
	}

	matches := hangarRegex.FindStringSubmatch(url)
	if matches == nil {
		log.FatalMsg("invalid url passed to plugin.handler.hangar", "url", url)
	}
	slug := matches[1]
	return &PluginHandler{
		url:  url,
		slug: slug,
	}
}

func (h *PluginHandler) Name() string {
	return "hangar"
}

func (h *PluginHandler) String() string {
	return fmt.Sprintf("%s:%s (%q)", "hangar", h.slug, h.url)
}

func (h *PluginHandler) LatestVersion() {

}

func (h *PluginHandler) GetMeta() (internal.PluginMeta, error) {
	url := fmt.Sprintf("https://hangar.papermc.io/api/v1/projects/%s", h.slug)
	var info PluginInfoJSON
	_, err := net.Get(url, "Could not get plugin info, plugin="+h.slug, &info)
	if err != nil {
		return internal.PluginMeta{}, err
	}

	var sourceURL string
	var supportURL string
	var wikiURL string

	for _, link := range info.Settings.Links[0].Links {
		if sourceURL != "" && supportURL != "" && wikiURL != "" {
			break
		}

		switch link.Name {
		case "Source":
			sourceURL = link.URL
		case "Support":
			supportURL = link.URL
		case "Wiki":
			wikiURL = link.URL
		}
	}

	return internal.PluginMeta{
		Title:       info.Name,
		Description: info.Description,
		Loaders:     nil,
		Updated:     info.LastUpdated,
		Downloads:   info.Stats.Downloads,
		Source:      sourceURL,
		Support:     supportURL,
		Wiki:        wikiURL,
	}, nil
}

func (h *PluginHandler) GetVersionList() ([]internal.PluginVersion, error) {
	url := fmt.Sprintf("https://hangar.papermc.io/api/v1/projects/%s/versions", h.slug)
	var info PluginVersionsJSON
	_, err := net.Get(url, "Could not get plugin info, plugin="+h.slug, &info)
	if err != nil {
		return nil, err
	}

	versions := make([]internal.PluginVersion, len(info.Result))

	for i, version := range info.Result {
		var loaders []string
		files := make([]internal.PluginFile, len(version.Downloads))
		for loader, download := range version.Downloads {
			loader = strings.ToLower(loader)
			loaders = append(loaders, loader)
			files = append(files, internal.PluginFile{
				URL:     download.DownloadURL,
				Name:    download.FileInfo.Name,
				Size:    download.FileInfo.SizeBytes,
				Loaders: []string{loader},
			})
		}

		versions[i] = internal.PluginVersion{
			ID:          version.Name,
			Name:        version.Name,
			Loaders:     loaders,
			PublishedAt: version.CreatedAt,
			Files:       files,
		}
	}

	return versions, nil
}
