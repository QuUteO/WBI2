package service

import (
	"CRUD/internal/model"
	"CRUD/internal/repository"
	"errors"
	"strings"
	"time"
)

var ErrInvalidEvent = errors.New("invalid event")

type CalendarRepository struct {
	repository.Repository
}

type Service interface {
	CreateCalendarEvent(event model.Event) (model.Event, error)
	UpdateCalendarEvent(event model.Event) error
	DeleteCalendarEvent(id int) error
	GetEventsForDay(userID int, day time.Time) ([]model.Event, error)
	GetEventsForWeek(userId int, startTime time.Time) ([]model.Event, error)
	GetEventsForMonth(userId int, startTime time.Time) ([]model.Event, error)
}

func NewService(repo repository.Repository) *CalendarRepository {
	return &CalendarRepository{repo}
}

func IsBusinessError(err error) bool {
	return errors.Is(err, repository.ErrEventNotFound)
}

func (c *CalendarRepository) CreateCalendarEvent(event model.Event) (model.Event, error) {
	if err := validateEventData(event); err != nil {
		return model.Event{}, err
	}

	return c.Repository.CreateEvent(event)
}

func (c *CalendarRepository) UpdateCalendarEvent(event model.Event) error {
	if event.ID <= 0 {
		return errors.New("invalid event id")
	}
	if err := validateEventData(event); err != nil {
		return err
	}

	return c.Repository.UpdateEvent(event)
}

func (c *CalendarRepository) DeleteCalendarEvent(id int) error {
	if id <= 0 {
		return errors.New("invalid event id")
	}

	return c.Repository.DeleteEvent(id)
}

func (c *CalendarRepository) GetEventsForDay(userID int, day time.Time) ([]model.Event, error) {
	return c.Repository.GetForPeriod(userID, day, day)
}

func (c *CalendarRepository) GetEventsForWeek(userId int, startTime time.Time) ([]model.Event, error) {
	endTime := startTime.AddDate(0, 0, 6)

	return c.Repository.GetForPeriod(userId, startTime, endTime)
}

func (c *CalendarRepository) GetEventsForMonth(userId int, startTime time.Time) ([]model.Event, error) {
	endTime := startTime.AddDate(0, 1, 0)

	return c.Repository.GetForPeriod(userId, startTime, endTime)
}

func validateEventData(event model.Event) error {
	if event.UserID <= 0 {
		return ErrInvalidEvent
	}
	if event.Date.IsZero() {
		return ErrInvalidEvent
	}
	if strings.TrimSpace(event.Text) == "" {
		return ErrInvalidEvent
	}

	return nil
}
