package handler

import (
	"CRUD/internal/repository"
	"CRUD/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateEventAcceptsJSONDate(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(
		http.MethodPost,
		"/create_event",
		strings.NewReader(`{"user_id":1,"date":"2023-12-31","event":"new year"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.CreateEvent(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), `"id":1`) {
		t.Fatalf("expected created event id in response, got %s", rec.Body.String())
	}
}

func TestEventsForDayInvalidDateReturnsBadRequest(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodGet, "/events_for_day?user_id=1&date=31-12-2023", nil)
	rec := httptest.NewRecorder()

	h.GetEventsForDay(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d with body %s", rec.Code, rec.Body.String())
	}
}

func TestDeleteMissingEventReturnsServiceUnavailable(t *testing.T) {
	h := newTestHandler()
	req := httptest.NewRequest(http.MethodPost, "/delete_event", strings.NewReader(`{"id":1}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.DeleteEvent(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d with body %s", rec.Code, rec.Body.String())
	}
}

func newTestHandler() *CalendarHandler {
	repo := repository.NewRepository()
	srv := service.NewService(repo)

	return NewHandler(srv)
}
