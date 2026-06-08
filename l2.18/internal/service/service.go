package service

import (
	"CRUD/internal/model"
	"CRUD/internal/repository"
	"errors"
	"time"
)

type CalendarRepository struct {
	repository.Repository
}

type Service interface {
	CreateCalendarEvent(event model.Event) error
	UpdateCalendarEvent(event model.Event) error
	DeleteCalendarEvent(id int) error
	GetEventsForDay(userID int, day time.Time) ([]model.Event, error)
	GetEventsForWeek(userId int, startTime time.Time) ([]model.Event, error)
	GetEventsForMonth(userId int, startTime time.Time) ([]model.Event, error)
}

func NewService(repo repository.Repository) *CalendarRepository {
	return &CalendarRepository{repo}
}

func (c *CalendarRepository) CreateCalendarEvent(event model.Event) error {
	if event.Text == "" {
		return errors.New("empty text")
	}

	return c.Repository.CreateEvent(event)
}

func (c *CalendarRepository) UpdateCalendarEvent(event model.Event) error {
	if event.Text == "" {
		return errors.New("event text cannot be empty")
	}
	if event.ID <= 0 {
		return errors.New("invalid event id")
	}

	return c.Repository.UpdateEvent(event)
}

func (c *CalendarRepository) DeleteCalendarEvent(id int) error {
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
