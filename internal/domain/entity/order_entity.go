package entity

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	Date       time.Time `json:"date" gorm:"not null"`
	TotalPrice float64   `json:"total_price" gorm:"not null"`
	PaymentID  uint      `json:"payment_id" gorm:"not null"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	Payment    *Payment  `json:"payment" gorm:"foreignKey:PaymentID"`
	Tickets    []Ticket  `json:"tickets" gorm:"foreignKey:OrderID"`
	gorm.Model
}

func (o *Order) TableName() string {
	return "orders"
}
