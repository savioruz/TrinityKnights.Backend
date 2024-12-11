package model

import "time"

type EventResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Time        time.Time `json:"time"`
	VenueID     uint      `json:"venue_id"`
	Status      string    `json:"status"`
}

type CreateEventRequest struct {
	Name        string    `json:"name" validate:"required,lte=100"`
	Description string    `json:"description" validate:"required,lte=255"`
	Date        time.Time `json:"date" validate:"required"`
	Time        time.Time `json:"time" validate:"required"`
	VenueID     uint      `json:"venue_id" validate:"required"`
}

type UpdateEventRequest struct {
	ID          uint      `param:"id" validate:"required"`
	Name        string    `json:"name" validate:"omitempty,lte=100"`
	Description string    `json:"description" validate:"omitempty,lte=255"`
	Date        time.Time `json:"date" validate:"omitempty"`
	Time        time.Time `json:"time" validate:"omitempty"`
	VenueID     uint      `json:"venue_id" validate:"omitempty"`
}

type GetEventRequest struct {
	ID uint `param:"id" validate:"required"`
}

type EventsRequest struct {
	Page  int    `query:"page" validate:"numeric"`
	Size  int    `query:"size" validate:"numeric"`
	Sort  string `query:"sort" validate:"omitempty,oneof=name date time venue_id"`
	Order string `query:"order" validate:"omitempty"`
}

type EventSearchRequest struct {
	Name        string `query:"name" validate:"omitempty,lte=100"`
	Description string `query:"description" validate:"omitempty,lte=255"`
	Date        string `query:"date" validate:"omitempty"`
	Time        string `query:"time" validate:"omitempty"`
	VenueID     uint   `query:"venue_id" validate:"omitempty"`
	Page        int    `query:"page" validate:"numeric"`
	Size        int    `query:"size" validate:"numeric"`
	Sort        string `query:"sort" validate:"omitempty,oneof=name date time"`
	Order       string `query:"order" validate:"omitempty"`
}

type EventQueryOptions struct {
	ID          *uint
	Name        *string
	Description *string
	Date        *string
	Time        *string
	VenueID     *uint
	Page        int
	Size        int
	Sort        string
	Order       string
}
