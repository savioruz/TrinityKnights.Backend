package entity

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     string    `json:"user_id" gorm:"not null"`
	Date       time.Time `json:"date" gorm:"not null"`
	TotalPrice float64   `json:"total_price" gorm:"not null"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	Payment    *Payment  `json:"payment" gorm:"foreignKey:OrderID"`
	Tickets    []Ticket  `json:"tickets" gorm:"foreignKey:OrderID"`
	Payments   []Payment `json:"payments" gorm:"foreignKey:OrderID"`
	gorm.Model
}

func (o *Order) TableName() string {
	return "orders"
}
