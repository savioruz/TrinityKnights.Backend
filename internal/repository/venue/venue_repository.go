package venue

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/repository"
	"gorm.io/gorm"
)

type VenueRepository interface {
	repository.Repository[entity.Venue]
	GetByID(db *gorm.DB, venue *entity.Venue, id uint) error
	GetPaginated(db *gorm.DB, venues *[]entity.Venue, opts model.VenueQueryOptions) (int64, error)
}
