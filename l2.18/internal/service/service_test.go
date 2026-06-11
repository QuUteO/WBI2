package service

import (
	"CRUD/internal/model"
	"CRUD/internal/repository"
	"testing"
	"time"
)

func TestCreateAndGetEvents(t *testing.T) {
	srv := NewService(repository.NewRepository())

	first := createEvent(t, srv, model.Event{
		UserID: 1,
		Date:   mustDate(t, "2023-12-31"),
		Text:   "first event",
	})
	second := createEvent(t, srv, model.Event{
		UserID: 1,
		Date:   mustDate(t, "2024-01-02"),
		Text:   "second event",
	})
	createEvent(t, srv, model.Event{
		UserID: 2,
		Date:   mustDate(t, "2023-12-31"),
		Text:   "another user event",
	})

	if first.ID == 0 || second.ID == 0 || first.ID == second.ID {
		t.Fatalf("expected generated unique ids, got %d and %d", first.ID, second.ID)
	}

	dayEvents, err := srv.GetEventsForDay(1, mustDate(t, "2023-12-31"))
	if err != nil {
		t.Fatalf("GetEventsForDay returned error: %v", err)
	}
	if len(dayEvents) != 1 || dayEvents[0].Text != "first event" {
		t.Fatalf("expected one event for day, got %#v", dayEvents)
	}

	weekEvents, err := srv.GetEventsForWeek(1, mustDate(t, "2023-12-31"))
	if err != nil {
		t.Fatalf("GetEventsForWeek returned error: %v", err)
	}
	if len(weekEvents) != 2 {
		t.Fatalf("expected two events for week, got %#v", weekEvents)
	}
}

func TestUpdateAndDeleteEvent(t *testing.T) {
	srv := NewService(repository.NewRepository())

	event := createEvent(t, srv, model.Event{
		UserID: 1,
		Date:   mustDate(t, "2023-12-31"),
		Text:   "old text",
	})

	event.Text = "new text"
	if err := srv.UpdateCalendarEvent(event); err != nil {
		t.Fatalf("UpdateCalendarEvent returned error: %v", err)
	}

	events, err := srv.GetEventsForDay(1, mustDate(t, "2023-12-31"))
	if err != nil {
		t.Fatalf("GetEventsForDay returned error: %v", err)
	}
	if len(events) != 1 || events[0].Text != "new text" {
		t.Fatalf("expected updated event, got %#v", events)
	}

	if err := srv.DeleteCalendarEvent(event.ID); err != nil {
		t.Fatalf("DeleteCalendarEvent returned error: %v", err)
	}
	if err := srv.DeleteCalendarEvent(event.ID); !IsBusinessError(err) {
		t.Fatalf("expected business error for missing event, got %v", err)
	}
}

func TestCreateInvalidEvent(t *testing.T) {
	srv := NewService(repository.NewRepository())

	if _, err := srv.CreateCalendarEvent(model.Event{
		UserID: 1,
		Date:   mustDate(t, "2023-12-31"),
		Text:   "",
	}); err == nil {
		t.Fatal("expected error for empty event text")
	}
}

func createEvent(t *testing.T, srv Service, event model.Event) model.Event {
	t.Helper()

	createdEvent, err := srv.CreateCalendarEvent(event)
	if err != nil {
		t.Fatalf("CreateCalendarEvent returned error: %v", err)
	}

	return createdEvent
}

func mustDate(t *testing.T, value string) time.Time {
	t.Helper()

	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		t.Fatalf("failed to parse date %q: %v", value, err)
	}

	return date
}
