package repository

import (
	"CRUD/internal/model"
	"fmt"
	"sync"
	"time"
)

// Repository описывает контракт (интерфейс)
type Repository interface {
	CreateEvent(event model.Event) error
	UpdateEvent(event model.Event) error
	DeleteEvent(id int) error
	GetForPeriod(userID int, start, end time.Time) ([]model.Event, error)
}

type RepositoryEvent struct {
	events map[int]model.Event
	mu     sync.RWMutex
	nextID int
}

func NewRepository() *RepositoryEvent {
	return &RepositoryEvent{
		events: make(map[int]model.Event),
		nextID: 1,
	}
}

func (r *RepositoryEvent) CreateEvent(event model.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	event.ID = r.nextID
	r.nextID++

	r.events[event.ID] = event
	return nil
}

func (r *RepositoryEvent) UpdateEvent(event model.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.events[event.ID]; !ok {
		return fmt.Errorf("event with id %d does not exist", event.ID)
	}

	r.events[event.ID] = event
	return nil
}

func (r *RepositoryEvent) DeleteEvent(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.events[id]; !ok {
		return fmt.Errorf("event with id %d does not exist", id)
	}

	delete(r.events, id)
	return nil
}

func (r *RepositoryEvent) GetForPeriod(userID int, start, end time.Time) ([]model.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.Event

	startTime := start.Truncate(24 * time.Hour)
	endTime := end.Truncate(24 * time.Hour).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	for _, ev := range r.events {
		if ev.UserID == userID {
			if (ev.Date.After(startTime) || ev.Date.Equal(startTime)) &&
				(ev.Date.Before(endTime) || ev.Date.Equal(endTime)) {
				result = append(result, ev)
			}
		}
	}

	return result, nil
}
