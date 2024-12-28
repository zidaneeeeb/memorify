package service

import (
	"context"
	"hbdtoyou/internal/template"
)

// CreateTemplate creates a new template and returns
// the created template ID.
func (s *service) CreateTemplate(ctx context.Context, reqTemplate template.Template) (string, error) {
	// validate fields
	err := validateTemplate(reqTemplate)
	if err != nil {
		return "", err
	}

	// update fields
	reqTemplate.CreateTime = s.timeNow()

	// get pg store client using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return "", err
	}

	// inserts template in pgstore
	templateID, err := pgStoreClient.CreateTemplate(ctx, reqTemplate)
	if err != nil {
		return "", err
	}

	return templateID, nil
}

// GetTemplateByID returns a template with the given
// template ID.
func (s *service) GetTemplateByID(ctx context.Context, templateID string) (template.Template, error) {
	// validate id
	if templateID == "" {
		return template.Template{}, template.ErrInvalidTemplateID
	}

	// get pg store client without transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return template.Template{}, err
	}

	// get template from pgstore
	result, err := pgStoreClient.GetTemplateByID(ctx, templateID)
	if err != nil {
		return template.Template{}, err
	}

	return result, nil
}

// GetTemplates returns all templates.
func (s *service) GetTemplates(ctx context.Context, filter template.GetTemplatesFilter) ([]template.Template, error) {
	// get pg store client without transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return nil, err
	}

	// get templates from pgstore
	result, err := pgStoreClient.GetTemplates(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateTemplate updates existing template
// with the given template data.
//
// UpdateTemplate do updates on all main attributes
// except ID, and CreateTime. So, make sure to
// use current values in the given data if do not want to
// update some specific attributes.
func (s *service) UpdateTemplate(ctx context.Context, reqTemplate template.Template) error {
	// validate id
	if reqTemplate.ID == "" {
		return template.ErrInvalidTemplateID
	}

	// validate fields
	err := validateTemplate(reqTemplate)
	if err != nil {
		return err
	}

	// update fields
	reqTemplate.UpdateTime = s.timeNow()

	// get pg store client using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return err
	}

	// inserts template in pgstore
	err = pgStoreClient.UpdateTemplate(ctx, reqTemplate)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTemplate delete a template
// with the given template id.
func (s *service) DeleteTemplateByID(ctx context.Context, templateID string) error {
	// validate id
	if templateID == "" {
		return template.ErrInvalidTemplateID
	}

	// get pg store client using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return err
	}

	// delete template in pgstore
	err = pgStoreClient.DeleteTemplateByID(ctx, templateID)
	if err != nil {
		return err
	}

	return nil
}

// validateTemplate validates fields of the given template
// whether its comply the predetermined rules.
func validateTemplate(reqTemplate template.Template) error {
	if reqTemplate.Name == "" {
		return template.ErrInvalidTemplateName
	}

	if _, valid := template.LabelList[reqTemplate.Label]; !valid {
		return template.ErrInvalidTemplateLabel
	}
	return nil
}
