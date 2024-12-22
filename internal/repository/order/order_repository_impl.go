package order

import (
	"github.com/TrinityKnights/Backend/internal/domain/entity"
	"github.com/TrinityKnights/Backend/internal/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *OrderRepositoryImpl) GetPaginatedOrders(db *gorm.DB, orders *[]entity.Order, page, size int, sort, order string) (int64, error) {
	var totalItems int64
	query := db.Model(&entity.Order{})

	// Add sorting
	validSortFields := map[string]bool{
		"ID":          true,
		"date":        true,
		"total_price": true,
		"created_at":  true,
	}

	validOrders := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if sort != "" && order != "" && validSortFields[sort] && validOrders[order] {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sort}, Desc: order == "desc"})
	} else {
		query = query.Order("created_at DESC")
	}

	// Get total count
	if err := query.Count(&totalItems).Error; err != nil {
		return 0, err
	}

	// Get paginated results
	offset := (page - 1) * size
	if err := query.Preload("Tickets").
		Preload("Tickets.Event").
		Preload("Payments").
		Offset(offset).
		Limit(size).
		Find(orders).Error; err != nil {
		return 0, err
	}

	return totalItems, nil
}
