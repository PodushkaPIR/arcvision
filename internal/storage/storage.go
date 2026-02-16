package storage

import (
	"encoding/json"
	"errors"
	"fatearcan/internal/domain"
	"fmt"
	"os"
	"slices"
)


type JSONStorage struct {
	cards []domain.TarotCard
	spreads []domain.Spread
}

func NewStorage(deckPath, spreadPath string) (*JSONStorage, error) {
	var storage JSONStorage

	if err := storage.loadCard(deckPath); err != nil {
		return nil, fmt.Errorf("failed to load cards: %w", err)
	}
	if err := storage.loadSpreads(spreadPath); err != nil {
		return nil, fmt.Errorf("failed to load spreads: %w", err)
	}

	return &storage, nil
}

var (
	ErrCardNotFound   = errors.New("card not found")
	ErrSpreadNotFound = errors.New("spread not found")
)

func (storage *JSONStorage) loadCard(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &storage.spreads)
}

func (storage *JSONStorage) loadSpreads(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &storage.spreads)
}

func (storage *JSONStorage) GetAllCards() []domain.TarotCard {
	return slices.Clone(storage.cards)
}

func (storage *JSONStorage) GetCardByID(id int) (*domain.TarotCard, error) {
	for _, card := range storage.cards {
		if card.ID == id {
			return &card, nil
		}
	}
	return nil, ErrCardNotFound
}

func (storage *JSONStorage) GetAllSpreads() []domain.Spread {
	return slices.Clone(storage.spreads)
}

func (storage *JSONStorage) GetSpreadByID(id int) (*domain.Spread, error) {
	for _, spread := range storage.spreads {
		if spread.ID == id {
			return &spread, nil
		}
	}
	return nil, ErrSpreadNotFound
}

