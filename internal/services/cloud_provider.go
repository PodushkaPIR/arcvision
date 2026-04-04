package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type CloudPredictor struct {
	log   *slog.Logger
	url   string
	key   string
	model string
}

func NewCloudPredictor(log *slog.Logger, url, key, model string) *CloudPredictor {
	return &CloudPredictor{
		log:   log,
		url:   url,
		key:   key,
		model: model,
	}
}

func (p *CloudPredictor) Generate(ctx context.Context, prompt string) (string, error) {
	p.log.Debug("Calling Cloud AI (OpenRouter)",
		slog.String("model", p.model),
		slog.Int("prompt_len", len(prompt)))

	payload := map[string]any{
		"model": p.model,
		"messages": []map[string]any{
			{"role": "system", "content": SystemPrompt},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
		"max_tokens":  2000,
	}

	data, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", p.url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+p.key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://arcvision.local")
	req.Header.Set("X-Title", "ArcVision")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	p.log.Debug("OpenRouter response",
		slog.Int("status", resp.StatusCode),
		slog.Int("body_len", len(body)))

	if resp.StatusCode != http.StatusOK {
		p.log.Error("OpenRouter API error",
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)))
		return "", fmt.Errorf("cloud AI status: %d, response: %s", resp.StatusCode, string(body))
	}

	if len(body) == 0 {
		p.log.Error("OpenRouter returned empty body")
		return "", fmt.Errorf("empty response from cloud AI")
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		p.log.Error("Failed to parse OpenRouter response",
			slog.String("error", err.Error()),
			slog.String("body", string(body)))
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Choices) == 0 {
		p.log.Error("OpenRouter returned no choices", slog.String("body", string(body)))
		return "", fmt.Errorf("empty choices from cloud AI")
	}

	p.log.Debug("AI response received", slog.Int("response_len", len(result.Choices[0].Message.Content)))
	return result.Choices[0].Message.Content, nil
}
