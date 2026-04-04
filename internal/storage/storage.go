package storage

import (
	"encoding/json"
	"errors"
	"fatearcan/internal/domain"
	"fmt"
	"math/rand/v2"
	"os"
)

type JSONStorage struct {
	cards   []domain.TarotCard
	spreads []domain.Spread
}

// NewStorage загружает данные при старте.
func NewStorage(deckPath, spreadPath string) (*JSONStorage, error) {
	s := &JSONStorage{} // Инициализируем указатель

	if err := s.loadFromFile(deckPath, &s.cards); err != nil {
		return nil, fmt.Errorf("loading deck: %w", err)
	}
	if err := s.loadFromFile(spreadPath, &s.spreads); err != nil {
		return nil, fmt.Errorf("loading spreads: %w", err)
	}

	return s, nil
}

// Generic метод для загрузки JSON
func (s *JSONStorage) loadFromFile(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// DataProvider - интерфейс для доступа к данным (Repository).
type DataProvider interface {
	GetSpread(id string) (domain.Spread, error)
	GetAllSpreads() []domain.Spread
	GetRandomCards(count int) ([]domain.TarotCard, error)
}

func (s *JSONStorage) GetAllSpreads() []domain.Spread {
	return s.spreads
}

func (s *JSONStorage) GetSpread(id string) (domain.Spread, error) {
	for _, spread := range s.spreads {
		if spread.ID == id {
			return spread, nil
		}
	}
	return domain.Spread{}, errors.New("spread not found")
}

func (s *JSONStorage) GetRandomCards(count int) ([]domain.TarotCard, error) {
	if count > len(s.cards) {
		return nil, errors.New("not enough cards in deck")
	}

	// Создаем копию слайса, чтобы не мешать порядок в оригинале
	deckCopy := make([]domain.TarotCard, len(s.cards))
	copy(deckCopy, s.cards)

	// Перемешиваем (Go 1.22 style)
	rand.Shuffle(len(deckCopy), func(i, j int) {
		deckCopy[i], deckCopy[j] = deckCopy[j], deckCopy[i]
	})

	return deckCopy[:count], nil
}
