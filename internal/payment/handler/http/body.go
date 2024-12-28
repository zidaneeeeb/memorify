package http

import (
	"hbdtoyou/internal/payment"
	"time"
)

// timeFormat denotes the standard time format used in
// payment HTTP handlers.
var timeFormat = "02/01/2006 3:04 PM -07:00"

type paymentHTTP struct {
	ID              *string `json:"id"`
	UserID          *string `json:"user_id"`
	UserName        *string `json:"user_name"`
	UserType        *string `json:"user_type"`
	UserQuota       *int    `json:"user_quota"`
	TemplateID      *string `json:"template_id"`
	TemplateName    *string `json:"template_name"`
	TemplateLabel   *string `json:"template_label"`
	ContentID       *string `json:"content_id"`
	Amount          *int    `json:"amount"`
	ProofPaymentURL *string `json:"proof_payment_url"`
	Date            *string `json:"date"`
	Status          *string `json:"status"`
}

func formatPayment(p payment.Payment) paymentHTTP {
	status := p.Status.String()
	userType := p.UserType.String()
	templateLabel := p.TemplateLabel.String()

	date := p.Date.Format(timeFormat)

	res := paymentHTTP{
		ID:              &p.ID,
		UserID:          &p.UserID,
		UserName:        &p.UserName,
		UserType:        &userType,
		UserQuota:       &p.UserQuota,
		TemplateID:      &p.TemplateID,
		TemplateName:    &p.TemplateName,
		TemplateLabel:   &templateLabel,
		ContentID:       &p.ContentID,
		Amount:          &p.Amount,
		ProofPaymentURL: &p.ProofPaymentURL,
		Date:            &date,
		Status:          &status,
	}

	return res
}

func (p paymentHTTP) parsePayment(out *payment.Payment) error {
	if p.ID != nil {
		out.ID = *p.ID
	}

	if p.UserID != nil {
		out.UserID = *p.UserID
	}

	if p.ContentID != nil {
		out.ContentID = *p.ContentID
	}

	if p.Amount != nil {
		out.Amount = *p.Amount
	}

	if p.ProofPaymentURL != nil {
		out.ProofPaymentURL = *p.ProofPaymentURL
	}

	if p.Date != nil && *p.Date != "" {
		endTime, err := time.Parse(timeFormat, *p.Date)
		if err != nil {
			return errInvalidPaymentDate
		}
		out.Date = endTime
	}

	if p.Status != nil {
		status, err := parsePaymentStatus(*p.Status)
		if err != nil {
			return err
		}
		out.Status = status
	}

	return nil
}

func parsePaymentStatus(req string) (payment.Status, error) {
	switch req {
	case payment.StatusDone.String():
		return payment.StatusDone, nil
	case payment.StatusPending.String():
		return payment.StatusPending, nil
	case payment.StatusRejected.String():
		return payment.StatusRejected, nil
	}

	return payment.StatusUnknown, errInvalidPaymentStatus
}
