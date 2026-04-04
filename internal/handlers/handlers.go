package handlers

import (
	"encoding/json"
	"fatearcan/internal/domain"
	"fatearcan/internal/services"
	"log/slog"
	"net/http"
)

type TarotHandler struct {
	service *services.TarotService
	log     *slog.Logger
}

func NewTarotHandler(service *services.TarotService, log *slog.Logger) *TarotHandler {
	return &TarotHandler{
		service: service,
		log:     log,
	}
}

func (h *TarotHandler) HandleSpreads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	spreads := h.service.GetSpreads()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(spreads); err != nil {
		h.log.Error("failed to encode spreads", slog.String("error", err.Error()))
	}
}

// Request payload
type readingRequest struct {
	SpreadID string `json:"spread_id"`
	Question string `json:"question"`
}

func (h *TarotHandler) HandleReading(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req readingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request", slog.String("error", err.Error()))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Вызываем бизнес-логику
	reading, err := h.service.CreateReading(r.Context(), req.SpreadID, req.Question)
	if err != nil {
		h.log.Error("failed to create reading", slog.String("error", err.Error()))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Отдаем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reading); err != nil {
		h.log.Error("failed to encode response", slog.String("error", err.Error()))
	}
}

type chatRequest struct {
	SystemPrompt string               `json:"system_prompt"`
	History      []domain.ChatMessage `json:"history"`
	SpreadID     string               `json:"spread_id"`
	Question     string               `json:"question"`
	Cards        []domain.DrawnCard   `json:"cards"`
}

type chatResponse struct {
	Message string `json:"message"`
}

func (h *TarotHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode chat request", slog.String("error", err.Error()))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.service.Chat(r.Context(), req.SystemPrompt, req.History, req.SpreadID, req.Question, req.Cards)
	if err != nil {
		h.log.Error("failed to process chat", slog.String("error", err.Error()))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatResponse{Message: response})
}
