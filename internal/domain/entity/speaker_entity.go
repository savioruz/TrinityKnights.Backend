package entity

import "gorm.io/gorm"

type Speaker struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string `json:"name" gorm:"not null"`
	Bio     string `json:"bio"`
	EventID uint   `json:"event_id" gorm:"not null"`
	Event   Event  `gorm:"foreignKey:EventID"`
	gorm.Model
}

func (s *Speaker) TableName() string {
	return "speakers"
}
