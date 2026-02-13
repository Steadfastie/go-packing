package domain

import "errors"

var (
	ErrInvalidAmount         = errors.New("amount must be greater than zero")
	ErrPackSizesNotConfigured = errors.New("pack sizes are not configured")
	ErrInvalidPackSizes      = errors.New("pack sizes must be non-empty unique positive integers")
	ErrCouldNotCalculate     = errors.New("could not calculate pack selection")
	ErrConcurrencyConflict    = errors.New("concurrency conflict")
)
