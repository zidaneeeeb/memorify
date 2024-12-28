package postgresql

import (
	"hbdtoyou/internal/auth"
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	ID         uuid.UUID  `db:"id"`
	Fullname   string     `db:"fullname"`
	Username   string     `db:"username"`
	Email      string     `db:"email"`
	Type       auth.Type  `db:"type"`
	Quota      int        `db:"quota"`
	CreateTime time.Time  `db:"create_time"`
	UpdateTime *time.Time `db:"update_time"`
}

func (dbData *UserModel) format() auth.User {
	u := auth.User{
		ID:         dbData.ID.String(),
		Fullname:   dbData.Fullname,
		Username:   dbData.Username,
		Email:      dbData.Email,
		Quota:      dbData.Quota,
		Type:       dbData.Type,
		CreateTime: dbData.CreateTime,
	}

	if dbData.UpdateTime != nil {
		u.UpdateTime = *dbData.UpdateTime
	}

	return u
}
