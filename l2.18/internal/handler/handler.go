package handler

import (
	"CRUD/internal/model"
	"CRUD/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type CalendarHandler struct {
	srv service.Service
}

func NewHandler(srv service.Service) *CalendarHandler {
	return &CalendarHandler{
		srv: srv,
	}
}

// Универсальный помощник для отправки успешных ответов (status 200)
func (h *CalendarHandler) sendSuccess(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"result": result,
	})
}

// Универсальный помощник для отправки ошибок
func (h *CalendarHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func (h *CalendarHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var event model.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
	}

	err := h.srv.CreateCalendarEvent(event)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
	}

	h.sendSuccess(w, nil)
}

func (h *CalendarHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var event model.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
	}

	err := h.srv.UpdateCalendarEvent(event)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
	}

	h.sendSuccess(w, nil)
}

func (h *CalendarHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	var event model.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
	}

	err := h.srv.DeleteCalendarEvent(event.ID)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
	}

	h.sendSuccess(w, nil)
}

func (h *CalendarHandler) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
	}

	date, err := time.Parse(r.URL.Query().Get("date"), "yyyy-mm-dd")
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
	}

	events, err := h.srv.GetEventsForDay(userID, date)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
	}

	h.sendSuccess(w, events)
}

func (h *CalendarHandler) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
	}

	date, err := time.Parse(r.URL.Query().Get("date"), "yyyy-mm-dd")
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
	}

	events, err := h.srv.GetEventsForWeek(userID, date)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
	}

	h.sendSuccess(w, events)
}

func (h *CalendarHandler) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
	}

	date, err := time.Parse(r.URL.Query().Get("date"), "yyyy-mm-dd")
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
	}

	events, err := h.srv.GetEventsForMonth(userID, date)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
	}

	h.sendSuccess(w, events)
}
