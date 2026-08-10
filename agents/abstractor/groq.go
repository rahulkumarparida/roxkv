package abstractor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GroqProvider implements the Provider interface for Groq using the OpenAI compatible API.
type GroqProvider struct {
	config ProviderConfig
	client *http.Client
}

// NewGroqProvider creates a new instance of GroqProvider.
func NewGroqProvider(config ProviderConfig) (*GroqProvider, error) {
	if config.Endpoint == "" {
		config.Endpoint = "https://api.groq.com/openai/v1"
	}
	timeout := time.Duration(config.Timeout) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &GroqProvider{
		config: config,
		client: &http.Client{Timeout: timeout},
	}, nil
}

// Chat sends a chat completion request to the Groq API.
// Name returns the provider identifier.
func (p *GroqProvider) Name() string {
	return "groq"
}

func (p *GroqProvider) Chat(ctx context.Context, messages []GenericMessage, tools []GenericToolDefinition, config ProviderConfig) (*ProviderResponse, error) {
	groqMsgs := make([]map[string]any, 0, len(messages))
	for _, msg := range messages {
		m := map[string]any{
			"role":    msg.Role,
			"content": msg.Content,
		}
		if msg.Role == "tool" {
			toolCallID := msg.ToolCallID
			if toolCallID == "" {
				toolCallID = "call_0"
			}
			m["tool_call_id"] = toolCallID
			if msg.ToolName != "" {
				m["name"] = msg.ToolName
			}
		}
		if len(msg.ToolCalls) > 0 {
			tcs := make([]map[string]any, 0, len(msg.ToolCalls))
			for _, tc := range msg.ToolCalls {
				argStr, _ := json.Marshal(tc.Arguments)
				tcs = append(tcs, map[string]any{
					"id":   tc.ID,
					"type": "function",
					"function": map[string]any{
						"name":      tc.Name,
						"arguments": string(argStr),
					},
				})
			}
			m["tool_calls"] = tcs
		}
		groqMsgs = append(groqMsgs, m)
	}

	reqBody := map[string]any{
		"model":       config.Model,
		"messages":    groqMsgs,
		"temperature": config.Temperature,
		"top_p":       config.TopP,
		"max_tokens":  config.MaxTokens,
	}

	if len(tools) > 0 {
		groqTools := make([]map[string]any, 0, len(tools))
		for _, t := range tools {
			groqTools = append(groqTools, buildOpenAIToolJSON(t))
		}
		reqBody["tools"] = groqTools
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := strings.TrimRight(config.Endpoint, "/")
	if endpoint == "" {
		endpoint = "https://api.groq.com/openai/v1"
	}
	endpoint = strings.TrimSuffix(endpoint, "/chat/completions")
	endpoint = strings.TrimSuffix(endpoint, "/v1")
	endpoint = strings.TrimRight(endpoint, "/")
	apiURL := fmt.Sprintf("%s/v1/chat/completions", endpoint)

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, NewProviderError("groq", 0, "failed to send request", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ClassifyHTTPError("groq", resp.StatusCode, string(body))
	}

	var groqResp struct {
		Choices []struct {
			Message struct {
				Content   string `json:"content"`
				ToolCalls []struct {
					ID       string `json:"id"`
					Function struct {
						Name      string `json:"name"`
						Arguments string `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &groqResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(groqResp.Choices) == 0 {
		return nil, ErrEmptyResponse
	}

	msg := groqResp.Choices[0].Message
	var responseToolCalls []GenericToolCall
	for _, tc := range msg.ToolCalls {
		var args map[string]any
		if tc.Function.Arguments != "" {
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
				return nil, NewProviderError("groq", 0, fmt.Sprintf("malformed tool call arguments JSON for %s: %v", tc.Function.Name, err), err)
			}
		}
		responseToolCalls = append(responseToolCalls, GenericToolCall{
			ID:        tc.ID,
			Name:      tc.Function.Name,
			Arguments: args,
		})
	}

	return &ProviderResponse{
		Text:      msg.Content,
		ToolCalls: responseToolCalls,
	}, nil
}
