package domain

type TarotCard struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Arcana      string `json:"arcana"`      // Major / Minor
	Suit        string `json:"suit"`        // Wands, etc. (пусто для старших)
	Description string `json:"description"` // Значение карты
	Image       string `json:"image"`
}

type Spread struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Positions   []string `json:"positions"` // Значения позиций (Прошлое, Настоящее...)
	Count       int      `json:"count"`     // Сколько карт тянуть
}

type Reading struct {
	Question       string      `json:"question"`
	Spread         Spread      `json:"spread"`
	Cards          []DrawnCard `json:"cards"`
	Interpretation string      `json:"interpretation"`
}

type DrawnCard struct {
	Card     TarotCard `json:"card"`
	Position string    `json:"position"`
}

type ChatMessage struct {
	Role    string `json:"role"`    // "user" or "assistant"
	Content string `json:"content"` // message text
}
