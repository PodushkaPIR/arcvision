package services

import (
	"context"
	"fatearcan/internal/domain"
	"fmt"
	"strings"
)

// DataProvider - интерфейс для доступа к данным (Repository).
// Определен ЗДЕСЬ, потому что он нужен ЗДЕСЬ.
type DataProvider interface {
	GetSpread(id string) (domain.Spread, error)
	GetAllSpreads() []domain.Spread
	GetRandomCards(count int) ([]domain.TarotCard, error)
}

// AIPredictor - интерфейс для нейросети.
type AIPredictor interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type TarotService struct {
	data DataProvider
	ai   AIPredictor
}

func NewTarotService(data DataProvider, ai AIPredictor) *TarotService {
	return &TarotService{
		data: data,
		ai:   ai,
	}
}

func (s *TarotService) GetSpreads() []domain.Spread {
	return s.data.GetAllSpreads()
}

func (s *TarotService) Chat(ctx context.Context, systemPrompt string, history []domain.ChatMessage, spreadID, question string, cards []domain.DrawnCard) (string, error) {
	var sb strings.Builder

	// Add system prompt
	sb.WriteString(systemPrompt)
	sb.WriteString("\n\n")

	// Add card context if available
	if len(cards) > 0 {
		sb.WriteString("Выпавшие карты:\n")
		for _, dc := range cards {
			sb.WriteString(fmt.Sprintf("- Позиция '%s': %s (%s). Значение: %s\n",
				dc.Position, dc.Card.Name, dc.Card.Arcana, dc.Card.Description))
		}
		sb.WriteString("\n")
	}

	// Add question context
	if question != "" {
		sb.WriteString(fmt.Sprintf("Вопрос кверента: %s\n\n", question))
	}

	// Add conversation history
	for _, msg := range history {
		sb.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
	}

	// Add current user message
	if question != "" {
		sb.WriteString(fmt.Sprintf("user: %s\n", question))
	}

	return s.ai.Generate(ctx, sb.String())
}

func (s *TarotService) CreateReading(ctx context.Context, spreadID string, question string) (*domain.Reading, error) {
	// 1. Получаем расклад
	spread, err := s.data.GetSpread(spreadID)
	if err != nil {
		return nil, fmt.Errorf("failed to get spread: %w", err)
	}

	// 2. Тянем карты
	cards, err := s.data.GetRandomCards(spread.Count)
	if err != nil {
		return nil, fmt.Errorf("failed to draw cards: %w", err)
	}

	// 3. Собираем структуру расклада
	var drawnCards []domain.DrawnCard
	for i, card := range cards {
		pos := "Доп. карта"
		if i < len(spread.Positions) {
			pos = spread.Positions[i]
		}
		drawnCards = append(drawnCards, domain.DrawnCard{
			Card:     card,
			Position: pos,
		})
	}

	// 4. Генерируем промпт для AI
	prompt := s.buildPrompt(question, spread.Name, drawnCards)

	// 5. Запрашиваем интерпретацию
	interpretation, err := s.ai.Generate(ctx, prompt)
	if err != nil {
		// Логируем ошибку, но можем вернуть расклад и без текста, если нужно
		return nil, fmt.Errorf("ai generation failed: %w", err)
	}

	return &domain.Reading{
		Question:       question,
		Spread:         spread,
		Cards:          drawnCards,
		Interpretation: interpretation,
	}, nil
}

func (s *TarotService) buildPrompt(question, spreadName string, cards []domain.DrawnCard) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Роль: Ты опытный таролог. Тон: мистический, но поддерживающий.\n"))
	sb.WriteString(fmt.Sprintf("Вопрос кверента: %s\n", question))
	sb.WriteString(fmt.Sprintf("Использованный расклад: %s\n\nВыпавшие карты:\n", spreadName))

	for _, dc := range cards {
		sb.WriteString(fmt.Sprintf("- Позиция '%s': %s (%s). Значение: %s\n",
			dc.Position, dc.Card.Name, dc.Card.Arcana, dc.Card.Description))
	}

	sb.WriteString("\nДай подробную интерпретацию расклада, связывая значения карт с позициями и вопросом.")
	return sb.String()
}
