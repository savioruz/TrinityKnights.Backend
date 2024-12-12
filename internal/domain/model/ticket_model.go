package model

// TicketResponse adalah model untuk merepresentasikan data tiket pada respons API.

// TicketResponse adalah model untuk merepresentasikan data tiket pada respons API.
type TicketResponse struct {
	ID         uint             `json:"id"`
	EventID    uint             `json:"event_id"`
	OrderID    uint             `json:"order_id"`
	Price      float64          `json:"price"`
	Type       string           `json:"type"`
	SeatNumber string           `json:"seat_number"`
	Event      *EventResponse   `json:"event,omitempty"`
	Order      *OrderResponse   `json:"order,omitempty"`
	Payment    *PaymentResponse `json:"payment,omitempty"`
}

// CreateTicketRequest adalah model untuk menerima data tiket saat pembuatan tiket.
type CreateTicketRequest struct {
	EventID    uint    `json:"event_id" validate:"required"`
	OrderID    uint    `json:"order_id" validate:"required"`
	Price      float64 `json:"price" validate:"required"`
	Type       string  `json:"type" validate:"required"`
	SeatNumber string  `json:"seat_number"`
}

// UpdateTicketRequest adalah model untuk menerima data tiket saat pembaruan tiket.
type UpdateTicketRequest struct {
	ID         uint    `param:"id" validate:"required"`
	EventID    uint    `json:"event_id" validate:"omitempty"`
	OrderID    uint    `json:"order_id" validate:"omitempty"`
	Price      float64 `json:"price" validate:"omitempty"`
	Type       string  `json:"type" validate:"omitempty"`
	SeatNumber string  `json:"seat_number" validate:"omitempty"`
}

// GetTicketRequest adalah model untuk mendapatkan data tiket berdasarkan ID.
type GetTicketRequest struct {
	ID uint `param:"id" validate:"required"`
}

// TicketsRequest adalah model untuk query tiket berdasarkan parameter pencarian.
type TicketsRequest struct {
	Page  int    `query:"page" validate:"numeric"`
	Size  int    `query:"size" validate:"numeric"`
	Sort  string `query:"sort" validate:"omitempty,oneof=event_id order_id price type seat_number"`
	Order string `query:"order" validate:"omitempty"`
}

// TicketSearchRequest adalah model untuk pencarian tiket berdasarkan query.
type TicketSearchRequest struct {
	EventID    uint   `query:"event_id" validate:"omitempty"`
	OrderID    uint   `query:"order_id" validate:"omitempty"`
	Price      string `query:"price" validate:"omitempty"`
	Type       string `query:"type" validate:"omitempty"`
	SeatNumber string `query:"seat_number" validate:"omitempty"`
	Page       int    `query:"page" validate:"numeric"`
	Size       int    `query:"size" validate:"numeric"`
	Sort       string `query:"sort" validate:"omitempty,oneof=event_id order_id price type seat_number"`
	Order      string `query:"order" validate:"omitempty"`
}

type TicketQueryOptions struct {
	ID         *uint
	EventID    *uint
	OrderID    *uint
	Price      *float64
	Type       *string
	SeatNumber *string
	Page       int
	Size       int
	Sort       string
	Order      string
}
