package template

import (
	"context"
	"time"
)

// Service is the interface for template service.
type Service interface {
	// CreateTemplate creates a new template and returns
	// the created template ID.
	CreateTemplate(ctx context.Context, reqTemplate Template) (string, error)

	// GetTemplateByID returns a template with the given
	// template ID.
	GetTemplateByID(ctx context.Context, templateID string) (Template, error)

	// GetTemplates returns all templates.
	GetTemplates(ctx context.Context, filter GetTemplatesFilter) ([]Template, error)

	// UpdateTemplate updates existing template
	// with the given template data.
	//
	// UpdateTemplate do updates on all main attributes
	// except ID, and CreateTime. So, make sure to
	// use current values in the given data if do not want to
	// update some specific attributes.
	UpdateTemplate(ctx context.Context, reqTemplate Template) error

	// DeleteTemplate delete a template
	// with the given template id.
	DeleteTemplateByID(ctx context.Context, templateID string) error
}

type GetTemplatesFilter struct {
	Label Label
}

type Template struct {
	ID           string
	Name         string
	Label        Label
	ThumbnailURI string
	CreateTime   time.Time
	UpdateTime   time.Time
}

// Label denotes the label of content.
type Label int

// Following constans are the known labels.
const (
	LabelUnknown Label = 0
	LabelFree    Label = 1
	LabelPremium Label = 2
)

// Following constans are the known label names.
var (
	// LabelList is a list of valid content label.
	LabelList = map[Label]struct{}{
		LabelFree:    {},
		LabelPremium: {},
	}

	// LabelName maps content label to it's string
	// representation.
	LabelName = map[Label]string{
		LabelFree:    "free",
		LabelPremium: "premium",
	}
)

// String implements the Stringer interface.
func (t Label) String() string {
	return LabelName[t]
}

// Value implements the Valuer interface.
func (t Label) Value() int {
	return int(t)
}
