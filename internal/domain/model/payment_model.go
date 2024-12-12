package model

// PaymentResponse adalah model untuk merepresentasikan data pembayaran pada respons API.
type PaymentResponse struct {
	ID            uint   `json:"id"`
	OrderID       uint   `json:"order_id"`
	Method        string `json:"method"`
	TransactionID string `json:"transaction_id"`
}

// CreatePaymentRequest adalah model untuk menerima data pembayaran saat pembuatan pembayaran.
type CreatePaymentRequest struct {
	OrderID       uint   `json:"order_id" validate:"required"`
	Method        string `json:"method" validate:"required"`
	TransactionID string `json:"transaction_id" validate:"required"`
}

// UpdatePaymentRequest adalah model untuk menerima data pembayaran saat pembaruan pembayaran.
type UpdatePaymentRequest struct {
	ID            uint   `param:"id" validate:"required"`
	OrderID       uint   `json:"order_id" validate:"omitempty"`
	Method        string `json:"method" validate:"omitempty"`
	TransactionID string `json:"transaction_id" validate:"omitempty"`
}

// GetPaymentRequest adalah model untuk menerima ID pembayaran saat meminta data pembayaran.
type GetPaymentRequest struct {
	ID uint `param:"id" validate:"required"`
}
