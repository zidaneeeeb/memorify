package postgresql

import (
	"hbdtoyou/internal/auth"
	"hbdtoyou/internal/payment"
	"hbdtoyou/internal/template"
	"time"
)

type paymentModel struct {
	ID              string         `db:"id"`
	UserID          string         `db:"user_id"`
	ContentID       string         `db:"content_id"`
	Amount          int            `db:"amount"`
	ProofPaymentURL string         `db:"proof_payment_url"`
	Date            time.Time      `db:"date"`
	UserName        string         `db:"user_name"`
	UserType        auth.Type      `db:"user_type"`
	UserQuota       int            `db:"user_quota"`
	TemplateID      string         `db:"template_id"`
	TemplateName    string         `db:"template_name"`
	TemplateLabel   template.Label `db:"template_label"`
	Status          payment.Status `db:"status"`
	CreateTime      time.Time      `db:"create_time"`
	UpdateTime      *time.Time     `db:"update_time"`
}

// format formats database struct into domain struct.
func (dbData *paymentModel) format() payment.Payment {
	p := payment.Payment{
		ID:              dbData.ID,
		UserID:          dbData.UserID,
		UserName:        dbData.UserName,
		UserType:        dbData.UserType,
		UserQuota:       dbData.UserQuota,
		TemplateID:      dbData.TemplateID,
		TemplateName:    dbData.TemplateName,
		TemplateLabel:   dbData.TemplateLabel,
		ContentID:       dbData.ContentID,
		Amount:          dbData.Amount,
		ProofPaymentURL: dbData.ProofPaymentURL,
		Date:            dbData.Date,
		Status:          dbData.Status,
		CreateTime:      dbData.CreateTime,
	}

	if dbData.UpdateTime != nil {
		p.UpdateTime = *dbData.UpdateTime
	}

	return p
}
