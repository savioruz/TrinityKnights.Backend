package model

type Response[T any] struct {
	Data   *T            `json:"data,omitempty"`
	Paging *PageMetadata `json:"paging,omitempty"`
	Error  *Error        `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PageMetadata struct {
	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

func NewResponse[T any](data T, paging *PageMetadata) *Response[T] {
	return &Response[T]{
		Data:   &data,
		Paging: paging,
	}
}

func NewErrorResponse[T any](code int, message string) *Response[T] {
	return &Response[T]{
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}
}
