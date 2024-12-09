package event

import (
	"context"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
)

type EventService interface {
	GetEventWithDetails(ctx context.Context, id uint) (*entity.Event, error)
}
