package event

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/repository"
	"gorm.io/gorm"
)

type EventRepository interface {
	repository.Repository[entity.Event]
	GetByID(db *gorm.DB, event *entity.Event, id uint) error
	GetPaginated(db *gorm.DB, events *[]entity.Event, opts *model.EventQueryOptions) (int64, error)
}
