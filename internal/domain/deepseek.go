package domain

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"text/template"
	"time"
)

var (
	deepSeekModel = "deepseek-chat"
	baseDeepSeek  = "https://api.deepseek.com/chat"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepSeekRequest struct {
	Model          string    `json:"model"`
	Messages       []Message `json:"messages"`
	Stream         bool      `json:"stream"`
	MaxTokens      int       `json:"max_tokens"`
	Temperature    float32   `json:"temperature"`
	ResponseFormat struct {
		Type string `json:"type"`
	} `json:"response_format"`
}

type DeepSeekResponse struct {
	Id      string `json:"id"`
	Choices []struct {
		FinishReason string  `json:"finish_reason"`
		Index        int     `json:"index"`
		Message      Message `json:"message"`
	} `json:"choices"`
	Created int    `json:"created"`
	Model   string `json:"model"`
}

type chatClient struct {
	client *http.Client
	key    string

	model     string
	maxTokens int
	temp      float32
	stream    bool

	responseFormat string
	systemPrompt   string
}

func NewDeepSeekClient(client *http.Client, apiKey, deepSeekModel, responseFormat string, maxTokens int, temp float32, stream bool, systemPrompt ...string) ChatBot {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	var sys string
	if len(systemPrompt) > 0 {
		for _, v := range systemPrompt {
			sys += v
		}
	}

	return &chatClient{
		client: client,
		key:    apiKey,

		model:     deepSeekModel,
		maxTokens: maxTokens,
		temp:      temp,
		stream:    stream,

		responseFormat: responseFormat,
		systemPrompt:   sys,
	}
}

const explainerSystemPrompt = `
You are a concise, accurate Japanese teacher and linguistic analyzer. 
For any provided Japanese sentence, produce a careful, step-by-step explanation including: kana reading, romaji, literal translation, natural translation(s), morpheme-by-morpheme gloss, grammatical analysis (each grammar point explained clearly with examples), nuance/register (politeness, offensiveness, formality), common learner mistakes, possible paraphrases/alternatives, and 2-5 practice exercises with answers. 
Use plain language, avoid speculation, and whenever a claim about usage or nuance is made, include a short justification (1 sentence).

Return output in JSON following the schema provided. Keep examples short and use only the words and structures relevant to the sentence unless you give a short contrast example. If the sentence contains offensive or sensitive language, flag it in the "nuance" field. Practice exercises should be clear and moderate to hard in difficulty.
`

func NewJapaneseExplainerClient(client *http.Client, apiKey string, maxTokens int) ChatBot {
	return NewDeepSeekClient(client, apiKey, deepSeekModel, "json_object", maxTokens, 0.1, false, explainerSystemPrompt)
}

func (c *chatClient) SetSystemPrompt(prompt string) {
	c.systemPrompt = prompt
}

func (c *chatClient) buildMessage(content string) []Message {
	var message []Message
	if c.systemPrompt != "" {
		message = append(message, Message{
			Role:    "system",
			Content: c.systemPrompt,
		})
	}

	message = append(message, Message{
		Role:    "user",
		Content: content,
	})

	return message
}

func (c *chatClient) request(content string) (io.Reader, error) {
	if len(content) == 0 {
		return nil, fmt.Errorf("text cannot be empty")
	}

	b, err := json.Marshal(DeepSeekRequest{
		Model:       c.model,
		Stream:      c.stream,
		MaxTokens:   c.maxTokens,
		Temperature: c.temp,
		Messages:    c.buildMessage(content),
		ResponseFormat: struct {
			Type string `json:"type"`
		}{
			Type: c.responseFormat,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	return bytes.NewReader(b), nil
}

func (c *chatClient) execute(ctx context.Context, body io.Reader, endpoint string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseDeepSeek+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.key)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("deepl returned status %d: %s", res.StatusCode, string(bodyBytes))
	}

	return res, nil
}

func (c *chatClient) handleResponse(res *http.Response) (string, error) {
	var deep DeepSeekResponse
	if err := json.NewDecoder(res.Body).Decode(&deep); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(deep.Choices) == 0 {
		return "", fmt.Errorf("no choices")
	}

	if deep.Choices[0].FinishReason != "stop" {
		return "", fmt.Errorf("no stop reason")
	}

	message := deep.Choices[0].Message.Content
	if message == "" {
		return "", fmt.Errorf("no message")
	}

	return message, nil
}

// The prompt is made by ChatGPT btw. :slightly_smiling_face:
var t = template.Must(template.New("base").Parse(`
Analyze this Japanese sentence and output JSON strictly matching the schema.

Sentence: "{{.Input}}"

Schema fields required:
{
  "original": string,
  "kana": string,
  "romaji": string,
  "literal_translation": string,
  "natural_translations": [string],
  "gloss_lines": {
     "surface": string,        // the original split into morphemes
     "reading": string,        // corresponding kana/romaji per morpheme
     "gloss": string           // short gloss per morpheme (EN)
  },
  "word_by_word": [ { "token": string, "reading": string, "pos": string, "meaning": string } ],
  "grammar_points": [ { "point": string, "explanation": string, "similar_examples": [string] } ],
  "nuance_and_register": string,
  "common_errors": [string],
  "paraphrases_and_alternatives": [string],
  "practice_exercises": [ { "task": string, "answer": string } ],
  "confidence": "low|medium|high"
}

Do not include any extra fields. Keep each explanation short (1–3 sentences). If you cannot analyze some item, set its value to null and explain briefly in its field.
`))

func (c *chatClient) buildExplainPrompt(input string) (string, error) {
	var buf bytes.Buffer
	if err := t.Execute(&buf, map[string]interface{}{
		"Input": input,
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}

type AskResponse struct {
	Original            string   `json:"original"`
	Kana                string   `json:"kana"`
	Romaji              string   `json:"romaji"`
	LiteralTranslation  string   `json:"literal_translation"`
	NaturalTranslations []string `json:"natural_translations"`
	GlossLines          struct {
		Surface string `json:"surface"`
		Reading string `json:"reading"`
		Gloss   string `json:"gloss"`
	} `json:"gloss_lines"`
	WordByWord []struct {
		Token   string `json:"token"`
		Reading string `json:"reading"`
		Pos     string `json:"pos"`
		Meaning string `json:"meaning"`
	} `json:"word_by_word"`
	GrammarPoints []struct {
		Point           string   `json:"point"`
		Explanation     string   `json:"explanation"`
		SimilarExamples []string `json:"similar_examples"`
	} `json:"grammar_points"`
	NuanceAndRegister          string   `json:"nuance_and_register"`
	CommonErrors               []string `json:"common_errors"`
	ParaphrasesAndAlternatives []string `json:"paraphrases_and_alternatives"`
	PracticeExercises          []struct {
		Task   string `json:"task"`
		Answer string `json:"answer"`
	} `json:"practice_exercises"`
	Confidence string `json:"confidence"`
}

func (c *chatClient) validateJapanese(content string) string {
	matches := regexp.MustCompile(`[\p{Hiragana}\p{Katakana}\p{Han}ー々、。「」『』？！]+`).FindAllString(content, -1)
	if len(matches) == 0 {
		return ""
	}

	valid := matches[0]
	for _, m := range matches {
		if len(m) > len(valid) {
			valid = m
		}
	}

	return valid
}

func (c *chatClient) Ask(ctx context.Context, content string) (*AskResponse, error) {
	japanese := c.validateJapanese(content)
	if japanese == "" {
		return nil, fmt.Errorf("invalid japanese sentence")
	}

	exp, err := c.buildExplainPrompt(japanese)
	if err != nil {
		return nil, err
	}

	body, err := c.request(exp)
	if err != nil {
		return nil, err
	}

	res, err := c.execute(ctx, body, "/completions")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	stringRes, err := c.handleResponse(res)
	if err != nil {
		return nil, err
	}

	var ask AskResponse
	if err = json.Unmarshal([]byte(stringRes), &ask); err != nil {
		return nil, err
	}

	return &ask, nil
}
