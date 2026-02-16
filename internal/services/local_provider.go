package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type LocalPredictor struct {
	log      *slog.Logger
	url      string
	model    string
}

func NewLocalPredictor(log *slog.Logger, url, model string) *LocalPredictor {
	return &LocalPredictor{log: log, url: url, model: model}
}

func (pred *LocalPredictor) Generate(ctx context.Context, prompt string) (string, error) {

	pred.log.Debug("Calling local ai api", slog.String("url", pred.url))

	payload := map[string]interface{} {
		"modelUri": fmt.Sprintf("gpt://%s/%s", pred.folderID, pred.model),
		"messages": []map[string]string {
			{"role": "user", "text": prompt},
		},
	}

	data, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", pred.url, bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Api-Key "+pred.key)


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

