package abstractor

import (
	"context"
	"testing"
)

type mockProvider struct {
	name          string
	receivedModel string
	receivedTemp  float64
	called        bool
	shouldFail    bool
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) Chat(ctx context.Context, messages []GenericMessage, tools []GenericToolDefinition, config ProviderConfig) (*ProviderResponse, error) {
	m.called = true
	m.receivedModel = config.Model
	m.receivedTemp = config.Temperature
	if m.shouldFail {
		return nil, NewProviderError(m.name, 500, "mock error", ErrProviderUnavailable)
	}
	return &ProviderResponse{Text: "mock success"}, nil
}

func TestExecuteWithFailover_Propagation(t *testing.T) {
	mock := &mockProvider{name: "ollama"}
	opts := map[string]any{
		"temperature": 0.7,
	}

	resp, active, err := ExecuteWithFailover(context.Background(), mock, "custom-model", opts, []GenericMessage{}, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if active != "ollama" {
		t.Errorf("expected active provider ollama, got %s", active)
	}
	if resp.Text != "mock success" {
		t.Errorf("expected mock success, got %s", resp.Text)
	}
	if mock.receivedModel != "custom-model" {
		t.Errorf("expected custom-model, got %s", mock.receivedModel)
	}
	if mock.receivedTemp != 0.7 {
		t.Errorf("expected temperature 0.7, got %f", mock.receivedTemp)
	}
}
