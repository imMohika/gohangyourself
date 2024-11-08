package net

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/imMohika/gohangyourself/log"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

func Request(url string, notOkErr string) *http.Response {
	log.Debug("Request URL:", url)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		log.Error(err, "an error occurred while making request", "url", url)
		return nil
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err, "an error occurred while sending request", "url", url)
		return nil
	}

	if response.StatusCode != http.StatusOK {
		log.Error(errors.New(notOkErr), "unexpected status code", "status", response.StatusCode, "url", url)
	}

	log.Debug("Done request", "status", response.StatusCode, "url", url)
	return response
}

func Get(url string, notOkErr string, content interface{}) (int, error) {
	resp := Request(url, notOkErr)

	log.Debug("decoding json")
	err := json.NewDecoder(resp.Body).Decode(&content)
	if err != nil {
		log.Error(err, "an error occurred while decoding json")
		return -1, err
	}

	err = resp.Body.Close()
	if err != nil {
		log.Error(err, "failed to close response body")
		return -1, err
	}

	return resp.StatusCode, nil
}
func GetGJSON(url string, notOkErr string) gjson.Result {
	resp := Request(url, notOkErr)

	log.Debug("decoding json")
	content, err := io.ReadAll(resp.Body)
	log.Error(err, "an error occurred while reading body")

	err = resp.Body.Close()
	log.Error(err, "failed to close response body")

	return gjson.ParseBytes(content)
}
