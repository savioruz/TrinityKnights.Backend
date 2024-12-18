package ticket

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TicketRepositoryImpl struct {
	repository.RepositoryImpl[entity.Ticket]
	Log *logrus.Logger
}

func NewTicketRepository(db *gorm.DB, log *logrus.Logger) *TicketRepositoryImpl {
	return &TicketRepositoryImpl{
		RepositoryImpl: repository.RepositoryImpl[entity.Ticket]{DB: db},
		Log:            log,
	}
}

func (r *TicketRepositoryImpl) CreateBatch(db *gorm.DB, tickets []*entity.Ticket) error {
	return db.Create(tickets).Error
}

func (r *TicketRepositoryImpl) Find(db *gorm.DB, opts *model.TicketQueryOptions) ([]*entity.Ticket, error) {
	query := db.Model(&entity.Ticket{})

	if opts.EventID != nil {
		query = query.Where("event_id = ?", *opts.EventID)
	}

	if len(opts.SeatNumbers) > 0 {
		query = query.Where("seat_number IN (?)", opts.SeatNumbers)
	}

	var tickets []*entity.Ticket
	if err := query.Find(&tickets).Error; err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *TicketRepositoryImpl) GetLastTicketNumber(db *gorm.DB, eventID uint, ticketType string) (*entity.Ticket, error) {
	var ticket entity.Ticket

	err := db.Transaction(func(tx *gorm.DB) error {
		// Lock the tickets table
		if err := tx.Exec("LOCK TABLE tickets IN SHARE ROW EXCLUSIVE MODE").Error; err != nil {
			return err
		}

		// Get the last ticket number for the specific event and type
		result := tx.Where("event_id = ? AND type = ?", eventID, ticketType).
			Order("CAST(SPLIT_PART(seat_number, '-', 2) AS INTEGER) DESC").
			First(&ticket)

		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return result.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}
