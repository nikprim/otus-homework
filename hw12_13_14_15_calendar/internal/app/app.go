package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nikprim/otus-homework/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage Storage
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, eventGUID uuid.UUID, event storage.Event) error
	DeleteEvent(ctx context.Context, eventGUID uuid.UUID) error
	FindEventsByInterval(ctx context.Context, startDateTime, endDateTime time.Time) ([]*storage.Event, error)
	FindEventByGUID(ctx context.Context, eventGUID uuid.UUID) (*storage.Event, error)
}

func New(storage Storage) *App {
	return &App{storage}
}
