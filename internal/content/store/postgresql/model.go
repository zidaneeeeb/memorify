package postgresql

import (
	"hbdtoyou/internal/content"
	"time"

	"github.com/google/uuid"
)

type contentModel struct {
	ID                    uuid.UUID      `db:"id"`
	UserID                string         `db:"user_id"`
	UserName              string         `db:"user_name"`
	TemplateID            uuid.UUID      `db:"template_id"`
	TemplateName          string         `db:"template_name"`
	TemplateLabel         string         `db:"template_label"`
	DetailContentJSONText string         `db:"detail_content_json_text"`
	Status                content.Status `db:"status"`
	CreateTime            time.Time      `db:"create_time"`
	UpdateTime            *time.Time     `db:"update_time"`
}

// format formats database struct into domain struct.
func (dbData *contentModel) format() content.Content {
	c := content.Content{
		ID:                    dbData.ID.String(),
		UserID:                dbData.UserID,
		UserName:              dbData.UserName,
		TemplateID:            dbData.TemplateID.String(),
		TemplateName:          dbData.TemplateName,
		TemplateLabel:         dbData.TemplateLabel,
		DetailContentJSONText: dbData.DetailContentJSONText,
		Status:                dbData.Status,
		CreateTime:            dbData.CreateTime,
	}

	if dbData.UpdateTime != nil {
		c.UpdateTime = *dbData.UpdateTime
	}

	return c
}
