package content

import "errors"

var (
	// ErrDataNotFound is returned when the wanted data is
	// not found.
	ErrDataNotFound = errors.New("data not found")

	// ErrInvalidUserID is returned when the given user ID
	// is invalid.
	ErrInvalidUserID = errors.New("invalid user id")

	// ErrInvalidContentID is returned when the given content
	// ID is invalid.
	ErrInvalidContentID = errors.New("invalid content id")

	// ErrInvalidContentType is returned when the given content
	// type is invalid.
	ErrInvalidContentType = errors.New("invalid content type")

	// ErrInvalidContentName is returned when the given content
	// name is invalid.
	ErrInvalidContentName = errors.New("invalid content name")

	// ErrInvalidTemplateID is returned when the given template
	// id is invalid.
	ErrInvalidTemplateID = errors.New("invalid template id")

	// ErrInvalidDetailContentJSONText is returned when the given
	// detail content json text is invalid.
	ErrInvalidDetailContentJSONText = errors.New("invalid detail content json text")

	// ErrInvalidContentStatus is returned when the given content
	// status is invalid.
	ErrInvalidContentStatus = errors.New("invalid content status")

	// ErrInvalidContentStatus is returned when the given condition
	// content access is invalid.
	ErrInvalidContentAccess = errors.New("invalid content access")
)
