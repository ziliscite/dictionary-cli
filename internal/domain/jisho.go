package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
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

const baseJisho = "https://jisho.org/api/v1/search/words"

func buildQuery(keyword string) string {
	params := make(url.Values)
	params.Add("keyword", keyword)
	return fmt.Sprintf("%s?%s", baseJisho, params.Encode())
}

type searcher struct {
	client *http.Client
}

func NewSearcher(client *http.Client) Searcher {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	return &searcher{
		client: client,
	}
}

func (s *searcher) getFromDictionary(keyword string) (io.ReadCloser, error) {
	get, err := s.client.Get(buildQuery(keyword))
	if err != nil {
		return nil, err
	}

	if get.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get error: %s", get.Status)
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

func (s *searcher) SearchRaw(keyword string) (*Jisho, error) {
	b, err := s.getFromDictionary(keyword)
	if err != nil {
		return nil, err
	}
	defer b.Close()

	return parseJisho(b)
}

func (s *searcher) Search(keyword string) ([]Information, error) {
	jisho, err := s.SearchRaw(keyword)
	if err != nil {
		return nil, err
	}

	return jisho.Data, nil
}
