package http

import (
	"hbdtoyou/internal/template"
	"net/http"
)

type templateHTTP struct {
	ID           *string `json:"id"`
	Name         *string `json:"name"`
	Label        *string `json:"label"`
	ThumbnailURI *string `json:"thumbnail_uri"`
}

func formatTemplate(t template.Template) templateHTTP {
	label := t.Label.String()

	return templateHTTP{
		ID:           &t.ID,
		Name:         &t.Name,
		Label:        &label,
		ThumbnailURI: &t.ThumbnailURI,
	}
}

func (t templateHTTP) parseTemplate(out *template.Template) error {
	if t.ID != nil {
		out.ID = *t.ID
	}

	if t.Name != nil {
		out.Name = *t.Name
	}

	if t.Label != nil {
		label, err := parseLabel(*t.Label)
		if err != nil {
			return err
		}

		out.Label = label
	}

	if t.ThumbnailURI != nil {
		out.ThumbnailURI = *t.ThumbnailURI
	}

	return nil
}

func parseLabel(req string) (template.Label, error) {
	switch req {
	case template.LabelFree.String():
		return template.LabelFree, nil
	case template.LabelPremium.String():
		return template.LabelPremium, nil
	}

	return template.LabelUnknown, errInvalidTemplateLabel
}

func (h *templatesHandler) parseHandleGetTemplatesQuery(r *http.Request) (template.GetTemplatesFilter, error) {
	query := r.URL.Query()

	res := template.GetTemplatesFilter{}

	labelParams := query.Get("label")
	if labelParams != "" {
		label, err := parseLabel(labelParams)
		if err != nil {
			return res, err
		}

		res.Label = label
	}

	return res, nil
}
