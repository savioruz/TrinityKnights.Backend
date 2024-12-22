package order

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/repository"
	"gorm.io/gorm"
)

type OrderRepository interface {
	repository.Repository[entity.Order]
	GetByID(db *gorm.DB, order *entity.Order, id uint) error
	GetByIDWithDetails(db *gorm.DB, order *entity.Order, id uint) error
	GetAllWithDetails(db *gorm.DB, orders *[]entity.Order) error
	GetPaginatedOrders(db *gorm.DB, orders *[]entity.Order, page, size int, sort, order string) (int64, error)
}
