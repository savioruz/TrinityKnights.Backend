package model

const (
	PaymentMethodBankTransfer = "bank_transfer"
	PaymentMethodQRIS         = "qris"
	// @TODO: Add other payment methods as needed
)

type PaymentRequest struct {
	OrderID       uint    `json:"order_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
}

type PaymentResponse struct {
	ID            uint    `json:"id"`
	OrderID       uint    `json:"order_id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"payment_method"`
	MidtransID    string  `json:"midtrans_id"`
	PaymentURL    string  `json:"payment_url"`
}

type MidtransCallbackRequest struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}
