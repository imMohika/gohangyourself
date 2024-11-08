package hangar

import "time"

type PluginInfoJSON struct {
	Name  string `json:"name"`
	Stats struct {
		Downloads int `json:"downloads"`
	} `json:"stats"`
	LastUpdated time.Time `json:"lastUpdated"`
	Description string    `json:"description"`
	Settings    struct {
		Links []struct {
			ID    int            `json:"id"`
			Type  string         `json:"type"`
			Title string         `json:"title"`
			Links []MetaLinkJSON `json:"links"`
		} `json:"links"`
	} `json:"settings"`
}

type MetaLinkJSON struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PluginVersionsJSON struct {
	Result []struct {
		CreatedAt time.Time `json:"createdAt"`
		Name      string    `json:"name"`
		Downloads map[string]struct {
			FileInfo struct {
				Name      string `json:"name"`
				SizeBytes int    `json:"sizeBytes"`
			} `json:"fileInfo"`
			DownloadURL string `json:"downloadUrl"`
		} `json:"downloads"`
	} `json:"result"`
}
