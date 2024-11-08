package internal

import "time"

type PluginMeta struct {
	Title       string
	Description string
	Loaders     []string
	Updated     time.Time
	Downloads   int
	Source      string
	Support     string
	Wiki        string
}

type PluginVersion struct {
	ID          string
	Name        string
	Loaders     []string
	PublishedAt time.Time
	Files       []PluginFile
}

type PluginFile struct {
	URL     string
	Name    string
	Size    int
	Loaders []string
}
