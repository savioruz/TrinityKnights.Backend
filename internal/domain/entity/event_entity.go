package entity

import (
	"time"

	"github.com/TrinityKnights/Backend/pkg/helper"
	"gorm.io/gorm"
)

type Event struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Date        time.Time      `json:"date" gorm:"type:date;not null"`
	Time        helper.SQLTime `json:"time" gorm:"type:time;not null"`
	VenueID     uint           `json:"venue_id" gorm:"not null"`
	Venue       Venue          `json:"venue" gorm:"foreignKey:VenueID"`
	gorm.Model
}

func (e *Event) TableName() string {
	return "events"
}
