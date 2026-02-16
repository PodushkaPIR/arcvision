package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type CloudPredictor struct {
	log      *slog.Logger
	url      string
	key      string
	folderID string
	model    string
}

func NewCloudPredictor(log *slog.Logger, url, key, folderID, model string) *CloudPredictor {
	return &CloudPredictor{log: log, url: url, key: key, folderID: folderID, model: model}
}

func (pred *CloudPredictor) Generate(ctx context.Context, prompt string) (string, error) {
	payload := map[string]interface{} {
		"modelUri": fmt.Sprintf("gpt://%s/%s", pred.folderID, pred.model),
		"messages": []map[string]string {
			{"role": "user", "text": prompt},
		},
	}

	data, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", pred.url, bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Api-Key "+pred.key)

	pred.log.Debug("Calling cloud ai api", slog.String("url", pred.url))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI API error: status %d", resp.StatusCode)
	}

	return "Базовый текст бебе", nil
}
