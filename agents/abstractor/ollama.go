package abstractor

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ollama/ollama/api"
)

// OllamaProvider implements the Provider interface for the Ollama API.
type OllamaProvider struct {
	client *api.Client
}

// NewOllamaProvider creates a new instance of OllamaProvider.
func NewOllamaProvider(config ProviderConfig) (*OllamaProvider, error) {
	parsedURL, err := url.Parse(config.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint URL: %w", err)
	}

	httpClient := http.DefaultClient
	if config.Timeout > 0 {
		httpClient = &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		}
	}

	client := api.NewClient(parsedURL, httpClient)

	return &OllamaProvider{
		client: client,
	}, nil
}

// Name returns the name of the provider.
func (p *OllamaProvider) Name() string {
	return "ollama"
}

// Chat sends a chat request to the Ollama API.
func (p *OllamaProvider) Chat(ctx context.Context, messages []GenericMessage, tools []GenericToolDefinition, config ProviderConfig) (*ProviderResponse, error) {
	apiMessages := convertMessagesToOllama(messages)
	apiTools := convertToolsToOllama(tools)

	options := map[string]interface{}{
		"temperature": config.Temperature,
		"top_p":       config.TopP,
	}
	if config.MaxTokens > 0 {
		options["num_predict"] = config.MaxTokens
	}

	stream := false
	req := &api.ChatRequest{
		Model:    config.Model,
		Messages: apiMessages,
		Tools:    apiTools,
		Options:  options,
		Stream:   &stream,
	}

	var responseText string
	var toolCalls []GenericToolCall

	respFunc := func(resp api.ChatResponse) error {
		responseText += resp.Message.Content
		for _, tc := range resp.Message.ToolCalls {
			toolCalls = append(toolCalls, GenericToolCall{
				ID:        "", // Ollama doesn't use call IDs
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments.ToMap(),
			})
		}
		return nil
	}

	err := p.client.Chat(ctx, req, respFunc)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return nil, fmt.Errorf("%w: %v", ErrConnectionRefused, err)
		}
		if errors.Is(err, context.DeadlineExceeded) || strings.Contains(err.Error(), "context deadline exceeded") {
			return nil, fmt.Errorf("%w: %v", ErrTimeout, err)
		}
		return nil, err
	}

	return &ProviderResponse{
		Text:      responseText,
		ToolCalls: toolCalls,
	}, nil
}

// convertMessagesToOllama converts generic messages to Ollama API messages.
func convertMessagesToOllama(messages []GenericMessage) []api.Message {
	var apiMessages []api.Message
	for _, m := range messages {
		msg := api.Message{
			Role:    m.Role,
			Content: m.Content,
		}
		if m.Role == "assistant" && len(m.ToolCalls) > 0 {
			for _, tc := range m.ToolCalls {
				args := api.NewToolCallFunctionArguments()
				for k, v := range tc.Arguments {
					args.Set(k, v)
				}
				msg.ToolCalls = append(msg.ToolCalls, api.ToolCall{
					Function: api.ToolCallFunction{
						Name:      tc.Name,
						Arguments: args,
					},
				})
			}
		}
		apiMessages = append(apiMessages, msg)
	}
	return apiMessages
}

// convertToolsToOllama converts generic tools to Ollama API tools.
func convertToolsToOllama(tools []GenericToolDefinition) []api.Tool {
	var apiTools []api.Tool
	for _, t := range tools {
		apiTool := api.Tool{
			Type: "function",
			Function: api.ToolFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters: api.ToolFunctionParameters{
					Type:     "object",
					Required: t.Parameters.Required,
				},
			},
		}

		properties := api.NewToolPropertiesMap()
		if properties != nil {
			for name, prop := range t.Parameters.Properties {
				apiProp := api.ToolProperty{
					Type:        api.PropertyType{prop.Type},
					Description: prop.Description,
				}
				if len(prop.Enum) > 0 {
					enumVals := make([]any, len(prop.Enum))
					for i, e := range prop.Enum {
						enumVals[i] = e
					}
					apiProp.Enum = enumVals
				}
				// Use type assertion since api.NewToolPropertiesMap is dynamic or returns an interface possibly,
				// or just use Set method as required by the instruction.
				properties.Set(name, apiProp)
			}
			apiTool.Function.Parameters.Properties = properties
		}

		apiTools = append(apiTools, apiTool)
	}
	return apiTools
}
