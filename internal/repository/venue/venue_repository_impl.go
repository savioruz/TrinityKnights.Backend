package venue

import (
	"strings"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type VenueRepositoryImpl struct {
	repository.RepositoryImpl[entity.Venue]
	Log *logrus.Logger
}

func NewVenueRepository(db *gorm.DB, log *logrus.Logger) *VenueRepositoryImpl {
	return &VenueRepositoryImpl{
		RepositoryImpl: repository.RepositoryImpl[entity.Venue]{DB: db},
		Log:            log,
	}
}

func (r *VenueRepositoryImpl) GetByID(db *gorm.DB, venue *entity.Venue, id uint) error {
	return db.Where("id = ?", id).Take(&venue).Error
}

func (r *VenueRepositoryImpl) GetPaginated(db *gorm.DB, venues *[]entity.Venue, opts *model.VenueQueryOptions) (int64, error) {
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.Size <= 0 {
		opts.Size = 10
	}

	query := r.buildPaginatedQuery(db, opts)

	// Get total count
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return 0, err
	}

	// Get paginated results
	offset := (opts.Page - 1) * opts.Size
	if err := query.Offset(offset).Limit(opts.Size).Find(venues).Error; err != nil {
		return 0, err
	}

	return totalCount, nil
}

func (r *VenueRepositoryImpl) buildPaginatedQuery(db *gorm.DB, opts *model.VenueQueryOptions) *gorm.DB {
	query := db.Model(&entity.Venue{})

	if opts.ID != nil && *opts.ID != 0 {
		query = query.Where("id = ?", *opts.ID)
	}
	if opts.Name != nil && *opts.Name != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+*opts.Name+"%")
	}
	if opts.Address != nil && *opts.Address != "" {
		query = query.Where("LOWER(address) LIKE LOWER(?)", "%"+*opts.Address+"%")
	}
	if opts.Capacity != nil && *opts.Capacity > 0 {
		query = query.Where("capacity = ?", *opts.Capacity)
	}
	if opts.City != nil && *opts.City != "" {
		query = query.Where("LOWER(city) LIKE LOWER(?)", "%"+*opts.City+"%")
	}
	if opts.State != nil && *opts.State != "" {
		query = query.Where("LOWER(state) LIKE LOWER(?)", "%"+*opts.State+"%")
	}
	if opts.Zip != nil && *opts.Zip != "" {
		query = query.Where("zip LIKE ?", "%"+*opts.Zip+"%")
	}

	// Add sorting with validation
	if opts.Sort != "" && opts.Order != "" {
		validSortFields := map[string]bool{
			"id":         true,
			"name":       true,
			"capacity":   true,
			"city":       true,
			"state":      true,
			"zip":        true,
			"created_at": true,
		}

		validOrders := map[string]bool{
			"asc":  true,
			"desc": true,
		}

		if validSortFields[strings.ToLower(opts.Sort)] && validOrders[strings.ToLower(opts.Order)] {
			orderClause := opts.Sort + " " + strings.ToLower(opts.Order)
			query = query.Order(orderClause)
		} else {
			query = query.Order("created_at DESC")
		}
	} else {
		query = query.Order("created_at DESC")
	}

	return query
}
