package domain

import "errors"

var (
	ErrInvalidDate = errors.New("invalid date")
	ErrInternal    = errors.New("internal error")

	ErrNotFound = errors.New("not found")
)
