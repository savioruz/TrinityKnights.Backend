package model

import "time"

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusPaid    PaymentStatus = "PAID"
	PaymentStatusExpired PaymentStatus = "EXPIRED"
)

type PaymentRequest struct {
	OrderID uint    `json:"order_id" validate:"required"`
	Amount  float64 `json:"amount" validate:"required"`
}

type PaymentResponse struct {
	ID         uint    `json:"id"`
	OrderID    uint    `json:"order_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	ExpiryDate string  `json:"expiry_date"`
	PaymentURL string  `json:"payment_url"`
}

type PaymentUpdateRequest struct {
	ID     uint          `json:"id" validate:"required"`
	Method string        `json:"method" validate:"required"`
	Status PaymentStatus `json:"status" validate:"required"`
}

type PaymentCallbackRequest struct {
	ID                 string    `json:"id" validate:"required"`
	ExternalID         string    `json:"external_id" validate:"required"`
	UserID             string    `json:"user_id" validate:"required"`
	IsHigh             bool      `json:"is_high" validate:"required"`
	PaymentMethod      string    `json:"payment_method" validate:"required"`
	Status             string    `json:"status" validate:"required"`
	MerchantName       string    `json:"merchant_name" validate:"required"`
	Amount             int       `json:"amount" validate:"required"`
	BankCode           *string   `json:"bank_code,omitempty" validate:"omitempty"`
	PaidAmount         int       `json:"paid_amount" validate:"required"`
	PaidAt             time.Time `json:"paid_at" validate:"required"`
	PayerEmail         string    `json:"payer_email" validate:"required"`
	Description        string    `json:"description" validate:"required"`
	Created            time.Time `json:"created" validate:"required"`
	Updated            time.Time `json:"updated" validate:"required"`
	Currency           string    `json:"currency" validate:"required"`
	PaymentChannel     string    `json:"payment_channel" validate:"required"`
	PaymentDestination string    `json:"payment_destination" validate:"required"`
}

type PaymentCallbackResponse struct {
	Status string `json:"status"`
}
