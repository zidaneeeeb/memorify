package template

import (
	"errors"
)

var (
	// ErrTemplateNotFound is returned when template is not found.
	ErrTemplateNotFound = errors.New("template not found")

	// ErrInvalidTemplateID is returned when template id is invalid.
	ErrInvalidTemplateID = errors.New("invalid template id")

	// ErrInvalidTemplateName is returned when template name is invalid.
	ErrInvalidTemplateName = errors.New("invalid template name")

	// ErrInvalidTemplateLabel is returned when template label is invalid.
	ErrInvalidTemplateLabel = errors.New("invalid template label")

	// ErrInvalidTemplateThumbnailURI is returned when template thumbnail uri is invalid.
	ErrInvalidTemplateThumbnailURI = errors.New("invalid template thumbnail uri")
)
