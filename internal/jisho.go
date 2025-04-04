package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Information struct {
	Slug     string   `json:"slug"`
	IsCommon bool     `json:"is_common"`
	JLPT     []string `json:"jlpt"`
	Japanese []struct {
		Word    string `json:"word,omitempty"`
		Reading string `json:"reading"`
	} `json:"japanese"`
	Senses []struct {
		EnglishDefinitions []string `json:"english_definitions"`
		PartsOfSpeech      []string `json:"parts_of_speech"`
	} `json:"senses"`
}

type Jisho struct {
	Data []Information `json:"data"`
}

const base = "https://jisho.org/api/v1/search/words"

func buildQuery(keyword string) string {
	params := make(url.Values)
	params.Add("keyword", keyword)
	return fmt.Sprintf("%s?%s", base, params.Encode())
}

var client = http.Client{}

func getFromDictionary(keyword string) (io.ReadCloser, error) {
	get, err := client.Get(buildQuery(keyword))
	if err != nil {
		return nil, err
	}

	if get.StatusCode != http.StatusOK {
		return nil, err
	}

	return get.Body, nil
}

func parseJisho(reader io.Reader) (*Jisho, error) {
	var data Jisho
	if err := json.NewDecoder(reader).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func Search(keyword string) (*Jisho, error) {
	b, err := getFromDictionary(keyword)
	if err != nil {
		return nil, err
	}

	return parseJisho(b)
}
