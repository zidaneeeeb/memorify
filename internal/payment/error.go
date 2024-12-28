package payment

import "errors"

var (
	// ErrUserAlreadyExist is returned when the given
	// user already exist based on the predefined
	// unique constraints.

	// ErrDataNotFound is returned when the desired data is
	// not found.
	ErrDataNotFound = errors.New("data not found")

	// ErrInvalidUserID is returned when the given user ID is
	// invalid.
	ErrInvalidUserID = errors.New("invalid user id")

	// ErrInvalidPaymentID is returned when the given payment ID is
	// invalid.
	ErrInvalidPaymentID = errors.New("invalid payment id")

	// ErrInvalidPaymentStatus is returned when the given payment status is
	// invalid.
	ErrInvalidPaymentStatus = errors.New("invalid payment status")

	// ErrInvalidAmount is returned when the given amount is
	// invalid.
	ErrInvalidAmount = errors.New("invalid amount")

	// ErrInvalidProofPaymentURL is returned when the given proof payment url is
	// invalid.
	ErrInvalidProofPaymentURL = errors.New("invalid proof payment url")

	// ErrInvalidPaymentDate is returned when the given payment date is
	// invalid.
	ErrInvalidPaymentDate = errors.New("invalid payment date")
)
