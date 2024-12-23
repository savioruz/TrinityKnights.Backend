// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphmodel

import (
	"time"

	"github.com/TrinityKnights/Backend/internal/domain/model"
)

type CreateTicketInput struct {
	EventID int     `json:"eventId"`
	Price   float64 `json:"price"`
	Type    string  `json:"type"`
	Count   int     `json:"count"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type EventsResponse struct {
	Data   []*model.EventResponse `json:"data,omitempty"`
	Paging *PageMetadata          `json:"paging,omitempty"`
	Error  *Error                 `json:"error,omitempty"`
}

type PageMetadata struct {
	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

type PaymentResponse struct {
	ID            int     `json:"id"`
	OrderID       int     `json:"orderId"`
	Amount        float64 `json:"amount"`
	TransactionID string  `json:"transactionId"`
	Method        *string `json:"method,omitempty"`
	Status        string  `json:"status"`
}

type PaymentsResponse struct {
	Data   []*PaymentResponse `json:"data,omitempty"`
	Paging *PageMetadata      `json:"paging,omitempty"`
	Error  *Error             `json:"error,omitempty"`
}

type Query struct {
}

type Response struct {
	Error  *Error        `json:"error,omitempty"`
	Paging *PageMetadata `json:"paging,omitempty"`
}

type TicketResponse struct {
	ID         string     `json:"id"`
	EventID    int        `json:"eventId"`
	OrderID    *int       `json:"orderId,omitempty"`
	Price      float64    `json:"price"`
	Type       string     `json:"type"`
	SeatNumber string     `json:"seatNumber"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

type TicketsResponse struct {
	Data   []*TicketResponse `json:"data,omitempty"`
	Paging *PageMetadata     `json:"paging,omitempty"`
	Error  *Error            `json:"error,omitempty"`
}

type UpdateEventInput struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Date        *string `json:"date,omitempty"`
	Time        *string `json:"time,omitempty"`
	VenueID     *int    `json:"venueId,omitempty"`
}

type UpdateTicketInput struct {
	EventID    *int     `json:"eventId,omitempty"`
	OrderID    *int     `json:"orderId,omitempty"`
	Price      *float64 `json:"price,omitempty"`
	Type       *string  `json:"type,omitempty"`
	SeatNumber *string  `json:"seatNumber,omitempty"`
}

type UpdateUserRequest struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

type UpdateVenueInput struct {
	Name     *string `json:"name,omitempty"`
	Address  *string `json:"address,omitempty"`
	Capacity *int    `json:"capacity,omitempty"`
	City     *string `json:"city,omitempty"`
	State    *string `json:"state,omitempty"`
	Zip      *string `json:"zip,omitempty"`
}

type VenuesResponse struct {
	Data   []*model.VenueResponse `json:"data,omitempty"`
	Paging *PageMetadata          `json:"paging,omitempty"`
	Error  *Error                 `json:"error,omitempty"`
}
