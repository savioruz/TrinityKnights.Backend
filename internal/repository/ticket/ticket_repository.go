package ticket

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/domain/model"
	"github.com/TrinityKnights/Backend/internal/repository"
	"gorm.io/gorm"
)

type TicketRepository interface {
	repository.Repository[entity.Ticket]
	CreateBatch(db *gorm.DB, tickets []*entity.Ticket) error
	Find(db *gorm.DB, filter *model.TicketQueryOptions) ([]*entity.Ticket, error)
	GetLastTicketNumber(db *gorm.DB, eventID uint, ticketType string) (*entity.Ticket, error)
}
