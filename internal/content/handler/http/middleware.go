package http

import (
	"context"
	"hbdtoyou/internal/auth"
	"log"
)

// checkAccessToken checks the given access token whether it
// is valid or not.
func checkAccessToken(ctx context.Context, auth auth.Service, token, userID, name string) error {
	tokenData, err := auth.ValidateToken(ctx, token)
	if err != nil {
		log.Printf("[Content HTTP][%s] Unauthorized error from ValidateToken. Err: %s\n", name, err.Error())
		return errUnauthorizedAccess
	}

	if userID != tokenData.UserID {
		return errInvalidUserID
	}

	return nil
}
