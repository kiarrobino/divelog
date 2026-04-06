package model

import "errors"

var (
	ErrDiveNotFound      = errors.New("dive not found")
	ErrInvalidDepth      = errors.New("depth must be greater than 0")
	ErrInvalidDuration   = errors.New("duration must be greater than 0")
	ErrInvalidRating     = errors.New("rating must be between 1 and 5")
	ErrInvalidDate       = errors.New("invalid date format, expected YYYY-MM-DD")
)
