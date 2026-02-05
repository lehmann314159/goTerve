package claude

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	apiURL     = "https://api.anthropic.com/v1/messages"
	apiVersion = "2023-06-01"
	model      = "claude-3-5-haiku-20241022"
)

// Client is a Claude API client
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// New creates a new Claude client
func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Message represents a message in the conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request represents an API request
type Request struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
}

// Response represents an API response
type Response struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// GenerateStory generates a Finnish story at the given CEFR level
func (c *Client) GenerateStory(cefrLevel, topic string) (story, translation string, err error) {
	if c.apiKey == "" {
		return "", "", fmt.Errorf("API key not configured")
	}

	prompt := buildStoryPrompt(cefrLevel, topic)

	req := Request{
		Model:     model,
		MaxTokens: 1024,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewReader(reqBody))
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", apiVersion)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp Response
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if apiResp.Error != nil {
		return "", "", fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	if len(apiResp.Content) == 0 {
		return "", "", fmt.Errorf("empty response from API")
	}

	// Parse the response to extract story and translation
	return parseStoryResponse(apiResp.Content[0].Text)
}

func buildStoryPrompt(cefrLevel, topic string) string {
	levelDescriptions := map[string]string{
		"A1": "very simple sentences, basic vocabulary (50-100 most common words), present tense only, short sentences",
		"A2": "simple sentences, common vocabulary, present and past tense, can use basic conjunctions",
		"B1": "moderate complexity, broader vocabulary, multiple tenses, compound sentences allowed",
		"B2": "more complex structures, idiomatic expressions allowed, varied vocabulary",
		"C1": "sophisticated language, nuanced expressions, complex grammar structures",
		"C2": "native-like fluency, literary expressions, full range of grammar",
	}

	levelDesc := levelDescriptions[cefrLevel]
	if levelDesc == "" {
		levelDesc = levelDescriptions["A1"]
	}

	return fmt.Sprintf(`Generate a short Finnish reading exercise for a %s level learner.

Topic: %s

Requirements:
- Write 4-6 sentences in Finnish appropriate for %s level: %s
- The story should be coherent and interesting
- Use proper Finnish grammar including correct cases and verb conjugations

Format your response EXACTLY like this:
FINNISH:
[Your Finnish story here]

ENGLISH:
[English translation here]

Do not include any other text or explanations.`, cefrLevel, topic, cefrLevel, levelDesc)
}

func parseStoryResponse(text string) (story, translation string, err error) {
	// Simple parsing - look for FINNISH: and ENGLISH: markers
	const finnishMarker = "FINNISH:"
	const englishMarker = "ENGLISH:"

	finnishIdx := bytes.Index([]byte(text), []byte(finnishMarker))
	englishIdx := bytes.Index([]byte(text), []byte(englishMarker))

	if finnishIdx == -1 || englishIdx == -1 {
		// Fallback: return the whole text as story
		return text, "", nil
	}

	story = text[finnishIdx+len(finnishMarker) : englishIdx]
	translation = text[englishIdx+len(englishMarker):]

	// Clean up whitespace
	story = trimWhitespace(story)
	translation = trimWhitespace(translation)

	return story, translation, nil
}

func trimWhitespace(s string) string {
	// Trim leading and trailing whitespace
	start := 0
	end := len(s)

	for start < end && (s[start] == ' ' || s[start] == '\n' || s[start] == '\r' || s[start] == '\t') {
		start++
	}

	for end > start && (s[end-1] == ' ' || s[end-1] == '\n' || s[end-1] == '\r' || s[end-1] == '\t') {
		end--
	}

	return s[start:end]
}