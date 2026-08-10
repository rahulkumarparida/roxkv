package abstractor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ClaudeProvider implements the Provider interface for Claude.
type ClaudeProvider struct {
	config ProviderConfig
	client *http.Client
}

// NewClaudeProvider creates a new ClaudeProvider.
func NewClaudeProvider(config ProviderConfig) (*ClaudeProvider, error) {
	if config.Endpoint == "" {
		config.Endpoint = "https://api.anthropic.com"
	}
	timeout := time.Duration(config.Timeout) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &ClaudeProvider{
		config: config,
		client: &http.Client{Timeout: timeout},
	}, nil
}

type claudeRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	System      string          `json:"system,omitempty"`
	Messages    []claudeMessage `json:"messages"`
	Tools       []claudeTool    `json:"tools,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	TopP        float64         `json:"top_p,omitempty"`
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"` // string or []claudeContentBlock
}

type claudeContentBlock struct {
	Type      string         `json:"type"`
	Text      string         `json:"text,omitempty"`
	ID        string         `json:"id,omitempty"`
	Name      string         `json:"name,omitempty"`
	Input     map[string]any `json:"input,omitempty"`
	ToolUseID string         `json:"tool_use_id,omitempty"`
	Content   string         `json:"content,omitempty"`
}

type claudeTool struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	InputSchema GenericToolParameters `json:"input_schema"`
}

type claudeResponse struct {
	Content []claudeContentBlock `json:"content"`
}

// Chat sends a chat completion request to Claude.
// Name returns the provider identifier.
func (p *ClaudeProvider) Name() string {
	return "claude"
}

func (p *ClaudeProvider) Chat(ctx context.Context, messages []GenericMessage, tools []GenericToolDefinition, config ProviderConfig) (*ProviderResponse, error) {
	if config.Endpoint == "" {
		config.Endpoint = "https://api.anthropic.com"
	}
	maxTokens := config.MaxTokens
	if maxTokens == 0 {
		maxTokens = 1024
	}

	reqBody := claudeRequest{
		Model:       config.Model,
		MaxTokens:   maxTokens,
		Temperature: config.Temperature,
		TopP:        config.TopP,
	}

	var systemPrompts []string
	for _, m := range messages {
		if m.Role == "system" {
			systemPrompts = append(systemPrompts, m.Content)
		} else {
			msg := claudeMessage{
				Role: m.Role,
			}
			if m.Role == "tool" {
				msg.Role = "user"
				toolUseID := m.ToolCallID
				if toolUseID == "" {
					toolUseID = "toolu_0"
				}
				msg.Content = []claudeContentBlock{
					{
						Type:      "tool_result",
						ToolUseID: toolUseID,
						Content:   m.Content,
					},
				}
			} else if len(m.ToolCalls) > 0 {
				var blocks []claudeContentBlock
				if m.Content != "" {
					blocks = append(blocks, claudeContentBlock{
						Type: "text",
						Text: m.Content,
					})
				}
				for _, tc := range m.ToolCalls {
					blocks = append(blocks, claudeContentBlock{
						Type:  "tool_use",
						ID:    tc.ID,
						Name:  tc.Name,
						Input: tc.Arguments,
					})
				}
				msg.Content = blocks
			} else {
				msg.Content = m.Content
			}
			reqBody.Messages = append(reqBody.Messages, msg)
		}
	}
	reqBody.System = strings.Join(systemPrompts, "\n")

	if len(tools) > 0 {
		for _, t := range tools {
			reqBody.Tools = append(reqBody.Tools, claudeTool{
				Name:        t.Name,
				Description: t.Description,
				InputSchema: t.Parameters,
			})
		}
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, NewProviderError("claude", 0, fmt.Sprintf("failed to marshal request: %v", err), err)
	}

	url := fmt.Sprintf("%s/v1/messages", config.Endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, NewProviderError("claude", 0, fmt.Sprintf("failed to create request: %v", err), err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, ErrProviderUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ClassifyHTTPError("claude", resp.StatusCode, "")
	}

	var claudeResp claudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&claudeResp); err != nil {
		return nil, NewProviderError("claude", 0, "failed to decode response", err)
	}

	providerResp := &ProviderResponse{}
	var textBlocks []string

	for _, block := range claudeResp.Content {
		if block.Type == "text" {
			textBlocks = append(textBlocks, block.Text)
		} else if block.Type == "tool_use" {
			providerResp.ToolCalls = append(providerResp.ToolCalls, GenericToolCall{
				ID:        block.ID,
				Name:      block.Name,
				Arguments: block.Input,
			})
		}
	}

	providerResp.Text = strings.Join(textBlocks, "\n")

	return providerResp, nil
}
