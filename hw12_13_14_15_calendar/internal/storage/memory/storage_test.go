package memorystorage_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nikprim/otus-homework/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/nikprim/otus-homework/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/require"
)

func TestStorage_CUD(t *testing.T) {
	eventStorage := memorystorage.New()
	ctx := context.Background()

	description := "description"
	notifyBefore := time.Hour

	event := storage.Event{
		GUID:         uuid.New(),
		Title:        "Test",
		StartAt:      time.Now(),
		EndAt:        time.Now().Add(time.Hour * 1),
		Description:  &description,
		UserGUID:     uuid.New(),
		NotifyBefore: &notifyBefore,
	}

	err := eventStorage.CreateEvent(ctx, event)
	require.NoError(t, err)

	resultEvent, err := eventStorage.FindEventByGUID(ctx, event.GUID)
	require.NoError(t, err)
	require.NotNil(t, resultEvent)
	require.Equal(t, &event, resultEvent)

	err = eventStorage.CreateEvent(ctx, event)
	require.ErrorIs(t, err, storage.ErrEventAlreadyExists)

	newEvent := storage.Event{
		GUID:     uuid.New(),
		StartAt:  time.Now().Add(time.Hour * 10),
		EndAt:    time.Now().Add(time.Hour * 15),
		UserGUID: uuid.New(),
	}

	err = eventStorage.UpdateEvent(ctx, event.GUID, newEvent)
	require.NoError(t, err)

	newEvent.GUID = event.GUID

	resultEvent, err = eventStorage.FindEventByGUID(ctx, event.GUID)
	require.NoError(t, err)
	require.NotNil(t, resultEvent)
	require.Equal(t, &newEvent, resultEvent)

	err = eventStorage.UpdateEvent(ctx, uuid.New(), newEvent)
	require.ErrorIs(t, err, storage.ErrEventNotFound)

	err = eventStorage.DeleteEvent(ctx, event.GUID)
	require.NoError(t, err)

	resultEvent, err = eventStorage.FindEventByGUID(ctx, event.GUID)
	require.NoError(t, err)
	require.Nil(t, resultEvent)

	err = eventStorage.DeleteEvent(ctx, event.GUID)
	require.NoError(t, err)
}
