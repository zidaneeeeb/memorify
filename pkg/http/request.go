package http

import (
	"errors"
	"net/http"
	"strings"
)

// Followings are the known error returned from HTTP request
// utility functions.
var (
	// ErrBearerTokenNotFound is returned when there is no
	// bearer token in the HTTP request header.
	ErrBearerTokenNotFound = errors.New("bearer token not found")

	// ErrSourceNotFound is returned when there is no source
	// in the HTTP request header.
	ErrSourceNotFound = errors.New("source not found")

	// ErrUserIDNotFound is returned when there is no user ID
	// in the HTTP request header.
	ErrUserIDNotFound = errors.New("user id not found")
)

// GetBearerTokenFromHeader returns token value stored in HTTP
// request header.
//
// Value is stored in standard header: Authorization.
func GetBearerTokenFromHeader(r *http.Request) (string, error) {
	value := r.Header.Get("Authorization")
	if value == "" {
		return "", ErrBearerTokenNotFound
	}

	parts := strings.Split(value, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", ErrBearerTokenNotFound
	}

	return parts[1], nil
}

// GetSourceFromHeader returns school ID value stored in HTTP
// request header.
//
// Value is stored in custom header: X-Source.
func GetSourceFromHeader(r *http.Request) (string, error) {
	source := r.Header.Get("X-Source")
	if source == "" {
		return "", ErrSourceNotFound
	}
	return source, nil
}

// GetUserIDFromHeader returns user ID value stored in HTTP
// request header.
//
// Value is stored in custom header: X-UserID.
func GetUserIDFromHeader(r *http.Request) (string, error) {
	userID := r.Header.Get("X-UserID")
	if userID == "" {
		return "", ErrUserIDNotFound
	}
	return userID, nil
}
