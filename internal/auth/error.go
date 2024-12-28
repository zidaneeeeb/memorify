package auth

import "errors"

var (
	// ErrUserAlreadyExist is returned when the given
	// user already exist based on the predefined
	// unique constraints.
	ErrUserAlreadyExist   = errors.New("user already exist")
	ErrDataNotFound       = errors.New("data not found")
	ErrInvalidEmail       = errors.New("invalid user email")
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrInvalidUserType    = errors.New("invalid user type")
	ErrInvalidTokenEmail  = errors.New("invalid user token email")

	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)
