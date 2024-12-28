package postgresql

import (
	"hbdtoyou/internal/template"
	"time"

	"github.com/google/uuid"
)

type templateModel struct {
	ID           uuid.UUID      `db:"id"`
	Name         string         `db:"name"`
	Label        template.Label `db:"label"`
	ThumbnailURI string         `db:"thumbnail_uri"`
	CreateTime   time.Time      `db:"create_time"`
	UpdateTime   *time.Time     `db:"update_time"`
}

// format formats database struct into domain struct.
func (dbData *templateModel) format() template.Template {
	t := template.Template{
		ID:           dbData.ID.String(),
		Name:         dbData.Name,
		Label:        dbData.Label,
		ThumbnailURI: dbData.ThumbnailURI,
		CreateTime:   dbData.CreateTime,
	}

	if dbData.UpdateTime != nil {
		t.UpdateTime = *dbData.UpdateTime
	}

	return t
}
