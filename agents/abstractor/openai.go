package abstractor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OpenAIProvider implements the Provider interface for OpenAI.
type OpenAIProvider struct {
	config ProviderConfig
	client *http.Client
}

// NewOpenAIProvider creates a new OpenAIProvider.
func NewOpenAIProvider(config ProviderConfig) (*OpenAIProvider, error) {
	if config.Endpoint == "" {
		config.Endpoint = "https://api.openai.com"
	}
	timeout := time.Duration(config.Timeout) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &OpenAIProvider{
		config: config,
		client: &http.Client{Timeout: timeout},
	}, nil
}

type openAIRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Tools       []openAITool    `json:"tools,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	TopP        float64         `json:"top_p,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Stream      bool            `json:"stream"`
}

type openAIMessage struct {
	Role       string           `json:"role"`
	Content    string           `json:"content"`
	ToolCalls  []openAIToolCall `json:"tool_calls,omitempty"`
	ToolCallID string           `json:"tool_call_id,omitempty"`
}

type openAIToolCall struct {
	ID       string             `json:"id"`
	Type     string             `json:"type"`
	Function openAIToolFunction `json:"function"`
}

type openAIToolFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type openAITool struct {
	Type     string         `json:"type"`
	Function openAIFunction `json:"function"`
}

type openAIFunction struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Parameters  GenericToolParameters `json:"parameters"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Role      string           `json:"role"`
			Content   string           `json:"content"`
			ToolCalls []openAIToolCall `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
}

// Name returns the provider identifier.
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// Chat sends a chat completion request to OpenAI.
func (p *OpenAIProvider) Chat(ctx context.Context, messages []GenericMessage, tools []GenericToolDefinition, config ProviderConfig) (*ProviderResponse, error) {
	reqBody := openAIRequest{
		Model:       p.config.Model,
		Temperature: p.config.Temperature,
		TopP:        p.config.TopP,
		MaxTokens:   p.config.MaxTokens,
		Stream:      p.config.Stream,
	}

	for _, m := range messages {
		msg := openAIMessage{
			Role:    m.Role,
			Content: m.Content,
		}
		if m.Role == "tool" {
			msg.ToolCallID = "call_default"
		}
		if len(m.ToolCalls) > 0 {
			for _, tc := range m.ToolCalls {
				argsBytes, _ := json.Marshal(tc.Arguments)
				msg.ToolCalls = append(msg.ToolCalls, openAIToolCall{
					ID:   tc.ID,
					Type: "function",
					Function: openAIToolFunction{
						Name:      tc.Name,
						Arguments: string(argsBytes),
					},
				})
			}
		}
		reqBody.Messages = append(reqBody.Messages, msg)
	}

	if len(tools) > 0 {
		for _, t := range tools {
			reqBody.Tools = append(reqBody.Tools, openAITool{
				Type: "function",
				Function: openAIFunction{
					Name:        t.Name,
					Description: t.Description,
					Parameters:  t.Parameters,
				},
			})
		}
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, NewProviderError("openai", 0, fmt.Sprintf("failed to marshal request: %v", err), err)
	}

	url := fmt.Sprintf("%s/v1/chat/completions", p.config.Endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, NewProviderError("openai", 0, fmt.Sprintf("failed to create request: %v", err), err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.APIKey))

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, ErrProviderUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ClassifyHTTPError("openai", resp.StatusCode, "")
	}

	var openAIResp openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return nil, NewProviderError("openai", 0, "failed to decode response", err)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, ErrEmptyResponse
	}

	choice := openAIResp.Choices[0].Message

	providerResp := &ProviderResponse{
		Text: choice.Content,
	}

	for _, tc := range choice.ToolCalls {
		var args map[string]any
		if tc.Function.Arguments != "" {
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
				// ignore invalid args
			}
		}
		providerResp.ToolCalls = append(providerResp.ToolCalls, GenericToolCall{
			ID:        tc.ID,
			Name:      tc.Function.Name,
			Arguments: args,
		})
	}

	return providerResp, nil
}
