package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type LocalPredictor struct {
	log   *slog.Logger
	url   string
	model string
}

func NewLocalPredictor(log *slog.Logger, url, model string) *LocalPredictor {
	return &LocalPredictor{
		log:   log,
		url:   url,   // Обычно http://localhost:11434/api/generate
		model: model, // Например, "llama3" или "mistral"
	}
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	System string `json:"system"`
}

// Структура ответа от Ollama
type ollamaResponse struct {
	Response string `json:"response"`
}

func (p *LocalPredictor) Generate(ctx context.Context, prompt string) (string, error) {
	p.log.Debug("Calling Local AI (Ollama)", slog.String("model", p.model))

	payload := ollamaRequest{
		Model:  p.model,
		Prompt: prompt,
		System: SystemPrompt,
		Stream: false,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", p.url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("local AI status: %d", resp.StatusCode)
	}

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Response, nil
}
