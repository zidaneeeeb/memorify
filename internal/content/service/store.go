package service

import (
	"context"
	"hbdtoyou/internal/content"
)

type PGStore interface {
	NewClient(useTx bool) (PGStoreClient, error)
}

type PGStoreClient interface {
	// Commit commits the transaction.
	Commit() error
	// Rollback aborts the transaction.
	Rollback() error

	// CreateContent creates a new content and returns
	// the created content ID.
	CreateContent(ctx context.Context, reqContent content.Content) (string, error)

	// GetContentByID returns a content with the given
	// content ID.
	GetContentByID(ctx context.Context, contentID string) (content.Content, error)

	// GetContents returns all contents.
	GetContents(ctx context.Context, filter content.GetContentsFilter) ([]content.Content, error)

	// UpdateContent updates existing content
	// with the given content data.
	//
	// UpdateContent do updates on all main attributes
	// except ID, and CreateTime. So, make sure to
	// use current values in the given data if do not want to
	// update some specific attributes.
	UpdateContent(ctx context.Context, reqContent content.Content) error

	// DeleteContent delete a content
	// with the given content id.
	DeleteContentByID(ctx context.Context, contentID string) error
}
