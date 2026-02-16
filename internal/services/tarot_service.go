package services

import (
	"crypto/rand"
	"math/big"
	"fatearcan/internal/domain"
)

type TarotService struct {
	cards []domain.TarotCard
	spreads []domain.Spread
}

func NewTarotService(cards []domain.TarotCard, spreads []domain.Spread) *TarotService {
	return &TarotService{
		cards:   cards,
		spreads: spreads,
	}
}

func (service *TarotService) GetSpread()

func (service *TarotService) CreateReading(spreadID int, question string) {

}

func (service *TarotService) drawRandomCard() (domain.TarotCard, error) {
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(service.cards))))
	if err != nil {
		return domain.TarotCard{}, err
	}
	return service.cards[index.Int64()], nil
}
