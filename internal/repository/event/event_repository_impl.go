package event

import (
	"errors"
	"fmt"
	"strings"

	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EventRepositoryImpl struct {
	repository.RepositoryImpl[entity.Event]
	Log *logrus.Logger
}

func NewEventRepository(db *gorm.DB, log *logrus.Logger) *EventRepositoryImpl {
	return &EventRepositoryImpl{
		RepositoryImpl: repository.RepositoryImpl[entity.Event]{DB: db},
		Log:            log,
	}
}

func (r *EventRepositoryImpl) GetByID(db *gorm.DB, event *entity.Event, id uint) error {
	result := db.Where("id = ?", id).First(event)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return result.Error
	}
	return nil
}

func (r *EventRepositoryImpl) GetPaginated(db *gorm.DB, events *[]entity.Event, opts *model.EventQueryOptions) (int64, error) {
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
	if err := query.Offset(offset).Limit(opts.Size).Find(events).Error; err != nil {
		return 0, err
	}

	return totalCount, nil
}

func (r *EventRepositoryImpl) buildPaginatedQuery(db *gorm.DB, opts *model.EventQueryOptions) *gorm.DB {
	query := db.Model(&entity.Event{})

	if opts.Name != nil && *opts.Name != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+*opts.Name+"%")
	}
	if opts.Description != nil && *opts.Description != "" {
		query = query.Where("LOWER(description) LIKE LOWER(?)", "%"+*opts.Description+"%")
	}
	if opts.Date != nil && *opts.Date != "" {
		query = query.Where("DATE(date) = ?", *opts.Date)
	}
	if opts.Time != nil && *opts.Time != "" {
		query = query.Where("time = ?::time", *opts.Time)
	}
	if opts.VenueID != nil && *opts.VenueID != 0 {
		query = query.Where("venue_id = ?", *opts.VenueID)
	}

	// Add sorting
	if opts.Sort != "" {
		direction := "ASC"
		if strings.EqualFold(opts.Order, "DESC") {
			direction = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", opts.Sort, direction))
	}

	return query
}
