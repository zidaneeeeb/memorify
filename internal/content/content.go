package content

import (
	"context"
	"time"
)

// Service is the interface for content service.
type Service interface {
	// CreateContent creates a new content and returns
	// the created content ID.
	CreateContent(ctx context.Context, reqContent Content) (string, error)

	// GetContentByID returns a content with the given
	// content ID.
	GetContentByID(ctx context.Context, contentID string) (Content, error)

	// GetContents returns all contents.
	GetContents(ctx context.Context, filter GetContentsFilter) ([]Content, error)

	// UpdateContent updates existing content
	// with the given content data.
	//
	// UpdateContent do updates on all main attributes
	// except ID, and CreateTime. So, make sure to
	// use current values in the given data if do not want to
	// update some specific attributes.
	UpdateContent(ctx context.Context, reqContent Content) error

	// DeleteContent delete a content
	// with the given content id.
	DeleteContentByID(ctx context.Context, contentID string) error
}

// Content denotes the content.
type Content struct {
	ID                    string
	UserID                string
	TemplateID            string
	DetailContentJSONText string
	Status                Status
	CreateTime            time.Time
	UpdateTime            time.Time

	// derived
	UserName      string
	TemplateName  string
	TemplateLabel string
}

// Status denotes status of a content.
type Status int

// Followings are the known content status.
const (
	StatusUnknown  Status = 0
	StatusActive   Status = 1
	StatusInactive Status = 2
)

var (
	// StatusList is a list of valid contnet status.
	StatusList = map[Status]struct{}{
		StatusActive:   {},
		StatusInactive: {},
	}

	// StatusName maps content status to it's string
	// representation.
	StatusName = map[Status]string{
		StatusActive:   "active",
		StatusInactive: "inactive",
	}
)

// Value returns int value of a content status.
func (s Status) String() string {
	return StatusName[s]
}

// String returns string representaion of a content status.
func (s Status) Value() int {
	return int(s)
}

type GetContentsFilter struct {
	UserID        string
	TemplateID    string
	TemplateLabel string
	Status        Status
}
