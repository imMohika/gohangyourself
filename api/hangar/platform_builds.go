package hangar

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"time"
)

type BuildListJSON struct {
	// ProjectID   string   `json:"project_id"`
	//ProjectName string   `json:"project_name"`
	//Version     string   `json:"version"`
	Builds []Builds `json:"builds"`
}

type Builds struct {
	Build    int       `json:"build"`
	Time     time.Time `json:"time"`
	Channel  string    `json:"channel"`
	Promoted bool      `json:"promoted"`
	// Changes   []Changes `json:"changes"`
	//Downloads Downloads `json:"downloads"`
}

// type Changes struct {
//	Commit  string `json:"commit"`
//	Summary string `json:"summary"`
//	Message string `json:"message"`
//}
//
//type Application struct {
//	Name   string `json:"name"`
//	Sha256 string `json:"sha256"`
//}
//
//type Downloads struct {
//	Application Application `json:"application"`
//}
//

func GetBuildList(platform string, version string) ([]Builds, error) {
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s/builds", platform, version)
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

	var buildList BuildListJSON
	err = json.Unmarshal(body, &buildList)
	if err != nil {
		return nil, err
	}

	builds := buildList.Builds
	slices.Reverse(builds)
	return builds, nil
}
