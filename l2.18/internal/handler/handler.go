package handler

import (
	"CRUD/internal/model"
	"CRUD/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "2006-01-02"

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
		h.sendError(w, "method not allowed", http.StatusBadRequest)
		return
	}

	event, err := parseEventRequest(r, false, true)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdEvent, err := h.srv.CreateCalendarEvent(event)
	if err != nil {
		h.sendServiceError(w, err)
		return
	}

	h.sendSuccess(w, createdEvent)
}

func (h *CalendarHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "method not allowed", http.StatusBadRequest)
		return
	}

	event, err := parseEventRequest(r, true, true)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.srv.UpdateCalendarEvent(event)
	if err != nil {
		h.sendServiceError(w, err)
		return
	}

	h.sendSuccess(w, "event updated")
}

func (h *CalendarHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "method not allowed", http.StatusBadRequest)
		return
	}

	event, err := parseEventRequest(r, true, false)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.srv.DeleteCalendarEvent(event.ID)
	if err != nil {
		h.sendServiceError(w, err)
		return
	}

	h.sendSuccess(w, "event deleted")
}

func (h *CalendarHandler) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "method not allowed", http.StatusBadRequest)
		return
	}

	userID, date, err := parseQueryParams(r)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.srv.GetEventsForDay(userID, date)
	if err != nil {
		h.sendServiceError(w, err)
		return
	}

	h.sendSuccess(w, events)
}

func (h *CalendarHandler) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "method not allowed", http.StatusBadRequest)
		return
	}

	userID, date, err := parseQueryParams(r)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.srv.GetEventsForWeek(userID, date)
	if err != nil {
		h.sendServiceError(w, err)
		return
	}

	h.sendSuccess(w, events)
}

func (h *CalendarHandler) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "method not allowed", http.StatusBadRequest)
		return
	}

	userID, date, err := parseQueryParams(r)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.srv.GetEventsForMonth(userID, date)
	if err != nil {
		h.sendServiceError(w, err)
		return
	}

	h.sendSuccess(w, events)
}

func (h *CalendarHandler) sendServiceError(w http.ResponseWriter, err error) {
	if service.IsBusinessError(err) {
		h.sendError(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	h.sendError(w, err.Error(), http.StatusBadRequest)
}

type eventRequest struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
	Text   string `json:"text"`
	Event  string `json:"event"`
}

func parseEventRequest(r *http.Request, requireID, requireDetails bool) (model.Event, error) {
	req, err := readEventRequest(r)
	if err != nil {
		return model.Event{}, err
	}
	if req.Text == "" {
		req.Text = req.Event
	}

	if requireID && req.ID <= 0 {
		return model.Event{}, fmt.Errorf("invalid id")
	}

	event := model.Event{ID: req.ID}
	if !requireDetails {
		return event, nil
	}

	if req.UserID <= 0 {
		return model.Event{}, fmt.Errorf("invalid user_id")
	}
	if strings.TrimSpace(req.Text) == "" {
		return model.Event{}, fmt.Errorf("empty event")
	}

	date, err := parseDate(req.Date)
	if err != nil {
		return model.Event{}, err
	}

	event.UserID = req.UserID
	event.Date = date
	event.Text = req.Text
	return event, nil
}

func readEventRequest(r *http.Request) (eventRequest, error) {
	if strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		return readEventForm(r)
	}

	var req eventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return eventRequest{}, err
	}
	return req, nil
}

func readEventForm(r *http.Request) (eventRequest, error) {
	if err := r.ParseForm(); err != nil {
		return eventRequest{}, err
	}

	id, err := parseOptionalInt(r.Form.Get("id"), "id")
	if err != nil {
		return eventRequest{}, err
	}
	userID, err := parseOptionalInt(r.Form.Get("user_id"), "user_id")
	if err != nil {
		return eventRequest{}, err
	}

	text := r.Form.Get("text")
	if text == "" {
		text = r.Form.Get("event")
	}

	return eventRequest{
		ID:     id,
		UserID: userID,
		Date:   r.Form.Get("date"),
		Text:   text,
	}, nil
}

func parseQueryParams(r *http.Request) (int, time.Time, error) {
	query := r.URL.Query()

	userID, err := parseRequiredInt(query.Get("user_id"), "user_id")
	if err != nil {
		return 0, time.Time{}, err
	}

	date, err := parseDate(query.Get("date"))
	if err != nil {
		return 0, time.Time{}, err
	}

	return userID, date, nil
}

func parseRequiredInt(value, name string) (int, error) {
	if value == "" {
		return 0, fmt.Errorf("%s is required", name)
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil || parsedValue <= 0 {
		return 0, fmt.Errorf("invalid %s", name)
	}

	return parsedValue, nil
}

func parseOptionalInt(value, name string) (int, error) {
	if value == "" {
		return 0, nil
	}

	return parseRequiredInt(value, name)
}

func parseDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("date is required")
	}

	date, err := time.Parse(dateLayout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date")
	}

	return date, nil
}
