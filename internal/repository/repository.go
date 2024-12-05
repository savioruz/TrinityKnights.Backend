package repository

import "gorm.io/gorm"

type Repository[T any] interface {
	Create(db *gorm.DB, entity *T) error
	Update(db *gorm.DB, entity *T) error
	Delete(db *gorm.DB, entity *T) error
}
