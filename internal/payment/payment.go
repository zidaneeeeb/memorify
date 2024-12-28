package payment

import (
	"context"
	"hbdtoyou/internal/auth"
	"hbdtoyou/internal/template"
	"time"
)

// Service is the interface for payment service.
type Service interface {
	// CreatePayment creates a new payment and returns
	// the created payment ID.
	CreatePayment(ctx context.Context, reqPayment Payment) (string, error)

	// GetPaymentByID returns a payment with the given
	// payment ID.
	GetPaymentByID(ctx context.Context, paymentID string) (Payment, error)

	// GetPayments returns all payments.
	GetPayments(ctx context.Context) ([]Payment, error)

	// UpdatePayment updates existing payment
	// with the given payment data.
	//
	// UpdatePayment do updates on all main attributes
	// except ID, and CreateTime. So, make sure to
	// use current values in the given data if do not want to
	// update some specific attributes.
	UpdatePayment(ctx context.Context, reqPayment Payment) error
}

// Payment denotes the payment.
type Payment struct {
	ID              string
	UserID          string
	ContentID       string
	Amount          int
	ProofPaymentURL string
	Date            time.Time
	Status          Status
	CreateTime      time.Time
	UpdateTime      time.Time

	// derived attributes
	UserName      string
	UserType      auth.Type
	UserQuota     int
	TemplateID    string
	TemplateName  string
	TemplateLabel template.Label
}

// Status denotes status of a payment.
type Status int

// Following constans are the known payment status.
const (
	StatusUnknown  Status = 0
	StatusDone     Status = 1
	StatusPending  Status = 2
	StatusRejected Status = 3
)

var (
	// StatusList is a list of valid payment status.
	StatusList = map[Status]struct{}{
		StatusDone:     {},
		StatusPending:  {},
		StatusRejected: {},
	}

	// StatusName maps payment status to it's string
	// representation.
	StatusName = map[Status]string{
		StatusDone:     "done",
		StatusPending:  "pending",
		StatusRejected: "rejected",
	}
)

// Value returns int value of a payment status.
func (s Status) String() string {
	return StatusName[s]
}

// String returns string representaion of a payment status.
func (s Status) Value() int {
	return int(s)
}
