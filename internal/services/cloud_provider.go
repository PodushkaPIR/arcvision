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
	log      *slog.Logger
	url      string
	key      string
	folderID string
	model    string
}

func NewCloudPredictor(log *slog.Logger, url, key, folderID, model string) *CloudPredictor {
	return &CloudPredictor{
		log:      log,
		url:      url, // https://llm.api.cloud.yandex.net/foundationModels/v1/completion
		key:      key,
		folderID: folderID,
		model:    model, // yandexgpt-lite или yandexgpt
	}
}

func (p *CloudPredictor) Generate(ctx context.Context, prompt string) (string, error) {
	p.log.Debug("Calling Cloud AI (Yandex)", slog.String("model", p.model))

	// Формат YandexGPT
	payload := map[string]any{
		"modelUri": fmt.Sprintf("gpt://%s/%s", p.folderID, p.model),
		"completionOptions": map[string]any{
			"stream": false,
			"temperature": 0.6,
			"maxTokens": "2000",
		},
		"messages": []map[string]string{
			{"role": "system", "text": "Ты мистический таролог."},
			{"role": "user", "text": prompt},
		},
	}

	data, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", p.url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	
	// Yandex требует Api-Key или IAM token
	req.Header.Set("Authorization", "Api-Key "+p.key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("cloud AI status: %d", resp.StatusCode)
	}

	// Парсим сложный ответ Яндекса
	var result struct {
		Result struct {
			Alternatives []struct {
				Message struct {
					Text string `json:"text"`
				} `json:"message"`
			} `json:"alternatives"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Result.Alternatives) > 0 {
		return result.Result.Alternatives[0].Message.Text, nil
	}

	return "", fmt.Errorf("empty response from cloud AI")
}
