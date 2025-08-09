package domain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type deepLRequest struct {
	Text       []string `json:"text"`
	TargetLang string   `json:"target_lang"`
}

type deepLResponse struct {
	Translations []Translation `json:"translations"`
}

type Translation struct {
	DetectedSourceLanguage string `json:"detected_source_language"`
	Text                   string `json:"text"`
}

type TargetLang int

const (
	TargetJapanese TargetLang = iota
	TargetEnglish
	TargetIndonesia
)

func (t TargetLang) Code() string {
	return [...]string{"JA", "EN", "ID"}[t]
}

func (t TargetLang) String() string {
	if t < TargetJapanese || t > TargetIndonesia {
		return "Unknown"
	}

	return [...]string{"Japanese", "English", "Indonesian"}[t]
}

type translator struct {
	client *http.Client

	key string
}

func NewDeepLClient(apiKey string, client *http.Client) Translator {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	return &translator{
		client: client,
		key:    apiKey,
	}
}

func (t *translator) request(lang TargetLang, texts ...string) (io.Reader, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("text cannot be empty")
	}

	b, err := json.Marshal(deepLRequest{
		Text:       texts,
		TargetLang: lang.Code(),
	})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	return bytes.NewReader(b), nil
}

var baseDeep = "https://api.deepl.com/v2"

func (t *translator) execute(ctx context.Context, body io.Reader, endpoint string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseDeep+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "DeepL-Auth-Key "+t.key)

	res, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("deepl returned status %d: %s", res.StatusCode, string(bodyBytes))
	}

	return res, nil
}

func (t *translator) Translate(ctx context.Context, lang TargetLang, texts ...string) ([]Translation, error) {
	body, err := t.request(lang, texts...)
	if err != nil {
		return nil, err
	}

	res, err := t.execute(ctx, body, "/translate")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var deep deepLResponse
	if err = json.NewDecoder(res.Body).Decode(&deep); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return deep.Translations, nil
}
