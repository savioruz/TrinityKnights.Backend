package model

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
	ID                 string  `json:"id"`
	ExternalID         string  `json:"external_id"`
	UserID             string  `json:"user_id"`
	IsHigh             bool    `json:"is_high"`
	PaymentMethod      string  `json:"payment_method"`
	Status             string  `json:"status"`
	MerchantName       string  `json:"merchant_name"`
	Amount             int     `json:"amount"`
	BankCode           *string `json:"bank_code,omitempty"`
	PaidAmount         int     `json:"paid_amount"`
	PaidAt             string  `json:"paid_at"`
	PayerEmail         string  `json:"payer_email"`
	Description        string  `json:"description"`
	Created            string  `json:"created"`
	Updated            string  `json:"updated"`
	Currency           string  `json:"currency"`
	PaymentChannel     string  `json:"payment_channel"`
	PaymentDestination string  `json:"payment_destination"`
}

type PaymentCallbackResponse struct {
	Status string `json:"status"`
}
