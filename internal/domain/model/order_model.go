package model

import (
	"time"
)

type OrderResponse struct {
	ID         uint             `json:"id"`
	UserID     uint             `json:"user_id"`
	Date       time.Time        `json:"date"`
	TotalPrice float64          `json:"total_price"`
	PaymentID  uint             `json:"payment_id"`
	User       UserResponse     `json:"user"`
	Payment    PaymentResponse  `json:"payment"`
	Tickets    []TicketResponse `json:"tickets"`
}

type CreateOrderRequest struct {
	UserID     uint    `json:"user_id" validate:"required"`
	TotalPrice float64 `json:"total_price" validate:"required"`
	PaymentID  uint    `json:"payment_id" validate:"required"`
	Tickets    []uint  `json:"tickets" validate:"required"` // Ticket IDs
}

type UpdateOrderRequest struct {
	ID         uint    `param:"id" validate:"required"`
	UserID     uint    `json:"user_id" validate:"omitempty"`
	TotalPrice float64 `json:"total_price" validate:"omitempty"`
	PaymentID  uint    `json:"payment_id" validate:"omitempty"`
	Tickets    []uint  `json:"tickets" validate:"omitempty"` // Updated Ticket IDs
}

type GetOrderRequest struct {
	ID uint `param:"id" validate:"required"`
}

type OrdersRequest struct {
	Page  int    `query:"page" validate:"numeric"`
	Size  int    `query:"size" validate:"numeric"`
	Sort  string `query:"sort" validate:"omitempty,oneof=id user_id total_price date"`
	Order string `query:"order" validate:"omitempty,oneof=asc,desc"`
}

type OrderSearchRequest struct {
	UserID     uint   `query:"user_id" validate:"omitempty"`
	Date       string `query:"date" validate:"omitempty"`
	TotalPrice string `query:"total_price" validate:"omitempty"`
	PaymentID  uint   `query:"payment_id" validate:"omitempty"`
	Page       int    `query:"page" validate:"numeric"`
	Size       int    `query:"size" validate:"numeric"`
	Sort       string `query:"sort" validate:"omitempty,oneof=id user_id total_price date"`
	Order      string `query:"order" validate:"omitempty,oneof=asc,desc"`
}

type OrderQueryOptions struct {
	ID         *uint
	UserID     *uint
	Date       *string
	TotalPrice *string
	PaymentID  *uint
	Page       int
	Size       int
	Sort       string
	Order      string
}
