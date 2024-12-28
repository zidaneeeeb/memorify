package service

import (
	"context"
	"hbdtoyou/internal/template"
)

type PGStore interface {
	NewClient(useTx bool) (PGStoreClient, error)
}

type PGStoreClient interface {
	// Commit commits the transaction.
	Commit() error
	// Rollback aborts the transaction.
	Rollback() error

	// CreateTemplate creates a new template and returns
	// the created template ID.
	CreateTemplate(ctx context.Context, reqTemplate template.Template) (string, error)

	// GetTemplateByID returns a template with the given
	// template ID.
	GetTemplateByID(ctx context.Context, templateID string) (template.Template, error)

	// GetTemplates returns all templates.
	GetTemplates(ctx context.Context, filter template.GetTemplatesFilter) ([]template.Template, error)

	// UpdateTemplate updates existing template
	// with the given template data.
	//
	// UpdateTemplate do updates on all main attributes
	// except ID, and CreateTime. So, make sure to
	// use current values in the given data if do not want to
	// update some specific attributes.
	UpdateTemplate(ctx context.Context, reqTemplate template.Template) error

	// DeleteTemplate delete a template
	// with the given template id.
	DeleteTemplateByID(ctx context.Context, templateID string) error
}
