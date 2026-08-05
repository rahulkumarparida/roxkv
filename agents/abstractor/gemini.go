package abstractor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GeminiProvider implements the Provider interface using raw HTTP calls to the Google Generative AI REST API.
type GeminiProvider struct {
	config ProviderConfig
	client *http.Client
}

// NewGeminiProvider creates a new instance of GeminiProvider.
func NewGeminiProvider(config ProviderConfig) (*GeminiProvider, error) {
	if config.Endpoint == "" {
		config.Endpoint = "https://generativelanguage.googleapis.com"
	}
	timeout := time.Duration(config.Timeout) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &GeminiProvider{
		config: config,
		client: &http.Client{Timeout: timeout},
	}, nil
}

// Chat sends a chat completion request to the Gemini API.
// Name returns the provider identifier.
func (p *GeminiProvider) Name() string {
	return "gemini"
}

func (p *GeminiProvider) Chat(ctx context.Context, messages []GenericMessage, tools []GenericToolDefinition, config ProviderConfig) (*ProviderResponse, error) {
	var systemInstruction *map[string]any
	var contents []map[string]any

	for _, msg := range messages {
		if msg.Role == "system" {
			systemInstruction = &map[string]any{
				"parts": []map[string]any{{"text": msg.Content}},
			}
			continue
		}

		parts := []map[string]any{}
		if msg.Content != "" {
			parts = append(parts, map[string]any{"text": msg.Content})
		}

		if msg.Role == "user" {
			contents = append(contents, map[string]any{
				"role":  "user",
				"parts": parts,
			})
		} else if msg.Role == "assistant" {
			for _, tc := range msg.ToolCalls {
				parts = append(parts, map[string]any{
					"functionCall": map[string]any{
						"name": tc.Name,
						"args": tc.Arguments,
					},
				})
			}
			contents = append(contents, map[string]any{
				"role":  "model",
				"parts": parts,
			})
		} else if msg.Role == "tool" {
			for _, tc := range msg.ToolCalls {
				parts = append(parts, map[string]any{
					"functionResponse": map[string]any{
						"name": tc.Name,
						"response": map[string]any{
							"content": msg.Content,
						},
					},
				})
			}
			contents = append(contents, map[string]any{
				"role":  "user",
				"parts": parts,
			})
		}
	}

	reqBody := map[string]any{
		"contents": contents,
		"generationConfig": map[string]any{
			"temperature":     p.config.Temperature,
			"topP":            p.config.TopP,
			"maxOutputTokens": p.config.MaxTokens,
		},
	}

	if systemInstruction != nil {
		reqBody["systemInstruction"] = *systemInstruction
	}

	if len(tools) > 0 {
		reqBody["tools"] = []map[string]any{buildGeminiToolJSON(tools)}
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", p.config.Endpoint, p.config.Model, p.config.APIKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, NewProviderError("gemini", 0, "failed to send request", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ClassifyHTTPError("gemini", resp.StatusCode, string(body))
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text         string `json:"text,omitempty"`
					FunctionCall struct {
						Name string         `json:"name"`
						Args map[string]any `json:"args"`
					} `json:"functionCall,omitempty"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 {
		return nil, ErrEmptyResponse
	}

	var responseText string
	var responseToolCalls []GenericToolCall

	for _, part := range geminiResp.Candidates[0].Content.Parts {
		if part.Text != "" {
			responseText += part.Text
		}
		if part.FunctionCall.Name != "" {
			responseToolCalls = append(responseToolCalls, GenericToolCall{
				ID:        fmt.Sprintf("call_%d", time.Now().UnixNano()),
				Name:      part.FunctionCall.Name,
				Arguments: part.FunctionCall.Args,
			})
		}
	}

	return &ProviderResponse{
		Text:      responseText,
		ToolCalls: responseToolCalls,
	}, nil
}
