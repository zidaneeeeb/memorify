package service

import (
	"context"
	"hbdtoyou/internal/auth"
)

type PGStore interface {
	NewClient(useTx bool) (PGStoreClient, error)
}

type PGStoreClient interface {
	// Commit commits the transaction.
	Commit() error
	// Rollback aborts the transaction.
	Rollback() error

	// CreateUser creates a new user and returns
	// the created user ID.
	CreateUser(ctx context.Context, user auth.User) (string, error)

	// GetUserAuth returns a user with the given
	// filter email or id.
	GetUserAuth(ctx context.Context, filter auth.GetUserAuthFilter) (auth.User, error)

	// UpdateUser updates existing user
	// with the given user data.
	//
	// UpdateUser do updates on all main attributes
	// except ID, and CreateTime. So, make sure to
	// use current values in the given data if do not want to
	// update some specific attributes.
	UpdateUser(ctx context.Context, reqUser auth.User) error
}
