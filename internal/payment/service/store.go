package service

import (
	"context"
	"hbdtoyou/internal/payment"
)

type PGStore interface {
	NewClient(useTx bool) (PGStoreClient, error)
}

type PGStoreClient interface {
	// Commit commits the transaction.
	Commit() error
	// Rollback aborts the transaction.
	Rollback() error

	// CreatePayment creates a new payment and returns
	// the created payment ID.
	CreatePayment(ctx context.Context, reqPayment payment.Payment) (string, error)

	// GetPaymentByID returns a payment with the given
	// payment ID.
	GetPaymentByID(ctx context.Context, paymentID string) (payment.Payment, error)

	// GetPayments returns all payments.
	GetPayments(ctx context.Context) ([]payment.Payment, error)

	// UpdatePayment updates existing payment
	// with the given payment data.
	//
	// UpdatePayment do updates on all main attributes
	// except ID, and CreateTime. So, make sure to
	// use current values in the given data if do not want to
	// update some specific attributes.
	UpdatePayment(ctx context.Context, reqPayment payment.Payment) error
}
