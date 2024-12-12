package order

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderRepositoryImpl struct {
	repository.RepositoryImpl[entity.Order]
	Log *logrus.Logger
}

func NewOrderRepository(db *gorm.DB, log *logrus.Logger) *OrderRepositoryImpl {
	return &OrderRepositoryImpl{
		RepositoryImpl: repository.RepositoryImpl[entity.Order]{DB: db},
		Log:            log,
	}
}

func (r *OrderRepositoryImpl) GetByID(db *gorm.DB, order *entity.Order, id uint) error {
	return db.Where("id = ?", id).Take(&order).Error
}

func (r *OrderRepositoryImpl) GetByIDWithDetails(db *gorm.DB, order *entity.Order, id uint) error {
	return db.Preload("Tickets").
		Preload("User").
		Preload("Payment").
		Where("orders.id = ?", id).
		Take(&order).Error
}

func (r *OrderRepositoryImpl) GetAllWithDetails(db *gorm.DB, orders *[]entity.Order) error {
	return db.Preload("Tickets").
		Preload("User").
		Preload("Payment").
		Find(&orders).Error
}
