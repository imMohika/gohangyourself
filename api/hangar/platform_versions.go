package hangar

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
)

type VersionListJSON struct {
	// ProjectID     string   `json:"project_id"`
	// ProjectName   string   `json:"project_name"`
	// VersionGroups []string `json:"version_groups"`
	Versions []string `json:"versions"`
}

func GetVersionList(platform string) ([]string, error) {
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s", platform)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("failed to close response body",
				"error", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, url %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var versionList VersionListJSON
	err = json.Unmarshal(body, &versionList)
	if err != nil {
		return nil, err
	}

	versions := versionList.Versions
	slices.Reverse(versions)
	return versions, nil
}
