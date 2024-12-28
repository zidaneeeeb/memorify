package service

import (
	"context"
	"hbdtoyou/internal/auth"
	"hbdtoyou/internal/content"
	"hbdtoyou/internal/template"

	"github.com/google/uuid"
)

// CreateContent creates a new content and returns
// the created content ID.
func (s *service) CreateContent(ctx context.Context, reqContent content.Content) (string, error) {
	// validate fields
	err := validateContent(reqContent)
	if err != nil {
		return "", err
	}

	currentTemplate, err := s.template.GetTemplateByID(ctx, reqContent.TemplateID)
	if err != nil {
		return "", err
	}

	if currentTemplate.Label == template.LabelPremium {
		user, err := s.user.GetUserByID(ctx, reqContent.UserID)
		if err != nil {
			return "", err
		}

		if user.Type == auth.TypeFree || user.Quota == 0 {
			return "", content.ErrInvalidContentAccess
		}
	}

	// update fields
	reqContent.CreateTime = s.timeNow()

	// get pg store client using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return "", err
	}

	// inserts content in pgstore
	contentID, err := pgStoreClient.CreateContent(ctx, reqContent)
	if err != nil {
		return "", err
	}

	return contentID, nil
}

// GetContentByID returns a content with the given
// content ID.
func (s *service) GetContentByID(ctx context.Context, contentID string) (content.Content, error) {
	// validate id
	if contentID == "" {
		return content.Content{}, content.ErrInvalidContentID
	}

	// get pg store client without transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return content.Content{}, err
	}

	// get content from pgstore
	result, err := pgStoreClient.GetContentByID(ctx, contentID)
	if err != nil {
		return content.Content{}, err
	}

	return result, nil
}

// GetContents returns all contents.
func (s *service) GetContents(ctx context.Context, filter content.GetContentsFilter) ([]content.Content, error) {
	// get pg store client without transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return nil, err
	}

	// get contents from pgstore
	result, err := pgStoreClient.GetContents(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateContent updates existing content
// with the given content data.
//
// UpdateContent do updates on all main attributes
// except ID, and CreateTime. So, make sure to
// use current values in the given data if do not want to
// update some specific attributes.
func (s *service) UpdateContent(ctx context.Context, reqContent content.Content) error {
	// validate id
	if reqContent.ID == "" {
		return content.ErrInvalidContentID
	}

	// validate fields
	err := validateContent(reqContent)
	if err != nil {
		return err
	}

	// update fields
	reqContent.UpdateTime = s.timeNow()

	// get pg store client using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return err
	}

	// inserts task in pgstore
	err = pgStoreClient.UpdateContent(ctx, reqContent)
	if err != nil {
		return err
	}

	return nil
}

// DeleteContent delete a content
// with the given content id.
func (s *service) DeleteContentByID(ctx context.Context, contentID string) error {
	// validate id
	if contentID == "" {
		return content.ErrInvalidContentID
	}

	// get pg store client using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return err
	}

	// delete content in pgstore
	err = pgStoreClient.DeleteContentByID(ctx, contentID)
	if err != nil {
		return err
	}

	return nil
}

// validateContent validates fields of the given content
// whether its comply the predetermined rules.
func validateContent(reqContent content.Content) error {
	if reqContent.TemplateID == "" {
		return content.ErrInvalidTemplateID
	}

	if _, valid := content.StatusList[reqContent.Status]; !valid {
		return content.ErrInvalidContentStatus
	}

	id, err := uuid.Parse(reqContent.UserID)
	if err != nil || id == uuid.Nil {
		return content.ErrInvalidUserID
	}

	if reqContent.DetailContentJSONText == "" {
		return content.ErrInvalidDetailContentJSONText
	}

	return nil
}
