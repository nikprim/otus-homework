package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEventAlreadyExists = errors.New("event already exists")
	ErrEventNotFound      = errors.New("event not found")
)

type Event struct {
	GUID         uuid.UUID
	Title        string
	StartAt      time.Time
	EndAt        time.Time
	Description  *string
	UserGUID     uuid.UUID
	NotifyBefore *time.Duration
}
