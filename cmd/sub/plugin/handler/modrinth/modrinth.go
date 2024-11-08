package modrinth

import (
	"fmt"
	"github.com/imMohika/gohangyourself/cmd/sub/plugin/internal"
	"github.com/imMohika/gohangyourself/log"
	"github.com/imMohika/gohangyourself/net"
	"regexp"
)

type PluginHandler struct {
	url  string
	slug string
}

var modrinthRegex = regexp.MustCompile(`^https:\/\/modrinth\.com\/plugin\/(.+)$`)

func FromURL(url string) *PluginHandler {
	matches := modrinthRegex.FindStringSubmatch(url)
	if matches == nil {
		log.FatalMsg("invalid url passed to plugin.handler.modrinth", "url", url)
	}
	slug := matches[1]
	return &PluginHandler{
		url:  url,
		slug: slug,
	}
}

func (h *PluginHandler) Name() string {
	return "modrinth"
}

func (h *PluginHandler) String() string {
	return fmt.Sprintf("%s:%s (%q)", "modrinth", h.slug, h.url)
}

func (h *PluginHandler) LatestVersion() {
	//TODO implement me
	panic("implement me")
}

func (h *PluginHandler) GetMeta() (internal.PluginMeta, error) {
	url := fmt.Sprintf("https://api.modrinth.com/v2/project/%s", h.slug)
	var info PluginInfoJSON
	_, err := net.Get(url, "Could not get plugin info, plugin="+h.slug, &info)
	if err != nil {
		return internal.PluginMeta{}, err
	}

	var supportURL string
	if info.DiscordURL != "" {
		supportURL = info.DiscordURL
	} else if info.IssuesURL != "" {
		supportURL = info.IssuesURL
	}

	return internal.PluginMeta{
		Title:       info.Title,
		Description: info.Description,
		Loaders:     info.Loaders,
		Updated:     info.Updated,
		Downloads:   info.Downloads,
		Source:      info.SourceURL,
		Support:     supportURL,
		Wiki:        info.WikiURL,
	}, nil
}

func (h *PluginHandler) GetVersionList() ([]internal.PluginVersion, error) {
	url := fmt.Sprintf("https://api.modrinth.com/v2/project/%s/version", h.slug)
	var info PluginVersionsJSON
	_, err := net.Get(url, "Could not get plugin info, plugin="+h.slug, &info)
	if err != nil {
		return nil, err
	}

	versions := make([]internal.PluginVersion, len(info))

	for i, version := range info {
		files := make([]internal.PluginFile, len(version.Files))
		for i, file := range version.Files {
			files[i] = internal.PluginFile{
				URL:     file.URL,
				Name:    file.Filename,
				Size:    file.Size,
				Loaders: version.Loaders,
			}
		}

		versions[i] = internal.PluginVersion{
			ID:          version.ID,
			Name:        version.Name,
			Loaders:     version.Loaders,
			PublishedAt: version.DatePublished,
			Files:       files,
		}
	}

	return versions, nil
}
