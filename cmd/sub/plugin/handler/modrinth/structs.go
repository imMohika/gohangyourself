package modrinth

import "time"

type PluginInfoJSON struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Updated     time.Time `json:"updated"`
	Downloads   int       `json:"downloads"`
	Loaders     []string  `json:"loaders"`
	IssuesURL   string    `json:"issues_url"`
	SourceURL   string    `json:"source_url"`
	WikiURL     string    `json:"wiki_url"`
	DiscordURL  string    `json:"discord_url"`
}

type PluginVersionsJSON []struct {
	Loaders       []string  `json:"loaders"`
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	DatePublished time.Time `json:"date_published"`
	Files         []struct {
		URL      string `json:"url"`
		Filename string `json:"filename"`
		Size     int    `json:"size"`
	} `json:"files"`
}
