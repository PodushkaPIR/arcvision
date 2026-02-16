package domain

import (
	"time"
)

type TarotCard struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Arcana      string `json:"arcana"`     // major/minor
	Suit        string `json:"suit,omitempty"` // wands, cups, swords, pentacles
	Value       string `json:"value,omitempty"` // ace, 2, 3... king
	Image       string `json:"image"`
}

type Spread struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Slots       []string `json:"slots"`
	MinCards    int      `json:"min_cards"`
	MaxCards    int      `json:"max_cards"`
}

type DrawnCard struct {
	Card       TarotCard `json:"card"`
	IsReversed bool      `json:"is_reversed"`
	Position   string    `json:"position"`
}

type Reading struct {
	ID        string      `json:"id"`
	Spread    Spread      `json:"spread"`
	Cards     []DrawnCard `json:"cards"`
	Question  string      `json:"question,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	Notes     string      `json:"notes,omitempty"`
}

type InterpretationRequest struct {
	Reading   Reading `json:"reading"`
	Question  string  `json:"question"`
	Language  string  `json:"language"` // ru/en
	Style     string  `json:"style"`   // brief/detailed
}
