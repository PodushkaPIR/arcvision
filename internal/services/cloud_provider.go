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
	p.log.Debug("Calling Cloud AI (OpenRouter)", slog.String("model", p.model))

	payload := map[string]any{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "user", "text": prompt},
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

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cloud AI status: %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("empty response from cloud AI")
}
