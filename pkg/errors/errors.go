package errors

import "errors"

var (
	ErrBadRequest         = errors.New("bad request")
	ErrNotFound           = errors.New("not found")
	ErrInternalServer     = errors.New("internal server error")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrValidation         = errors.New("validation error")
	ErrDuplicateEntry     = errors.New("duplicate entry")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
