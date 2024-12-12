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
