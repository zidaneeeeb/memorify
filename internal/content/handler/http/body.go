package http

import (
	"hbdtoyou/internal/content"
	"net/http"
)

type contentHTTP struct {
	ID                    *string `json:"id"`
	UserID                *string `json:"user_id"`
	Username              *string `json:"user_name"`
	TemplateID            *string `json:"template_id"`
	TemplateName          *string `json:"template_name"`
	TemplateLabel         *string `json:"template_label"`
	DetailContentJSONText *string `json:"detail_content_json_text"`
	Type                  *string `json:"type"`
	Status                *string `json:"status"`
}

func formatContent(c content.Content) contentHTTP {
	status := c.Status.String()

	return contentHTTP{
		ID:                    &c.ID,
		UserID:                &c.UserID,
		Username:              &c.UserName,
		TemplateID:            &c.TemplateID,
		TemplateName:          &c.TemplateName,
		TemplateLabel:         &c.TemplateLabel,
		Status:                &status,
		DetailContentJSONText: &c.DetailContentJSONText,
	}
}

func (c contentHTTP) parseContent(out *content.Content) error {
	if c.ID != nil {
		out.ID = *c.ID
	}

	if c.UserID != nil {
		out.UserID = *c.UserID
	}
	if c.TemplateID != nil {
		out.TemplateID = *c.TemplateID
	}

	if c.DetailContentJSONText != nil {
		out.DetailContentJSONText = *c.DetailContentJSONText
	}

	return nil
}

func parseContentStatus(req string) (content.Status, error) {
	switch req {
	case content.StatusActive.String():
		return content.StatusActive, nil
	case content.StatusInactive.String():
		return content.StatusInactive, nil
	}

	return content.StatusUnknown, errInvalidContentStatus
}

func (h *contentsHandler) parseHandleGetContentsQuery(r *http.Request) (content.GetContentsFilter, error) {
	query := r.URL.Query()

	res := content.GetContentsFilter{
		UserID:        query.Get("user_id"),
		TemplateID:    query.Get("template_id"),
		TemplateLabel: query.Get("template_label"),
	}

	statusParams := query.Get("status")
	if statusParams != "" {
		contentStatus, err := parseContentStatus(statusParams)
		if err != nil {
			return res, err
		}

		res.Status = contentStatus
	}
	return res, nil
}
