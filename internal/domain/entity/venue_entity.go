package entity

import "gorm.io/gorm"

type Venue struct {
	ID       uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string  `json:"name" gorm:"not null"`
	Address  string  `json:"address"`
	Capacity int     `json:"capacity"`
	City     string  `json:"city"`
	State    string  `json:"state"`
	Zip      string  `json:"zip" gorm:"column:zip_code"`
	Events   []Event `gorm:"foreignKey:VenueID"`
	gorm.Model
}

func (v *Venue) TableName() string {
	return "venues"
}
