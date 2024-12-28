package service

import (
	"context"
	"hbdtoyou/internal/auth"
	"hbdtoyou/internal/payment"
)

func (s *service) CreatePayment(ctx context.Context, reqPayment payment.Payment) (string, error) {
	// validate fields
	if reqPayment.ProofPaymentURL == "" {
		return "", payment.ErrInvalidProofPaymentURL
	}

	if reqPayment.Amount <= 0 {
		return "", payment.ErrInvalidAmount
	}

	// update fields
	reqPayment.CreateTime = s.timeNow()
	reqPayment.Status = payment.StatusPending

	// get pg store client using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return "", err
	}

	// inserts content in pgstore
	paymentID, err := pgStoreClient.CreatePayment(ctx, reqPayment)
	if err != nil {
		return "", err
	}

	user, err := s.user.GetUserByID(ctx, reqPayment.UserID)
	if err != nil {
		return "", err
	}

	user.Quota = 1
	user.Type = auth.TypePending

	err = s.user.UpdateUser(ctx, user)
	if err != nil {
		return "", err
	}

	return paymentID, nil
}

func (s *service) GetPaymentByID(ctx context.Context, paymentID string) (payment.Payment, error) {
	// validate id
	if paymentID == "" {
		return payment.Payment{}, payment.ErrInvalidPaymentID
	}

	// get pg store client without transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return payment.Payment{}, err
	}

	// get payment from pgstore
	result, err := pgStoreClient.GetPaymentByID(ctx, paymentID)
	if err != nil {
		return payment.Payment{}, err
	}

	return result, nil
}

func (s *service) GetPayments(ctx context.Context) ([]payment.Payment, error) {
	// get pg store client without transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return nil, err
	}

	// get payments from pgstore
	result, err := pgStoreClient.GetPayments(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) UpdatePayment(ctx context.Context, reqPayment payment.Payment) error {
	// validate id
	if reqPayment.ID == "" {
		return payment.ErrInvalidPaymentID
	}

	// update fields
	reqPayment.UpdateTime = s.timeNow()

	// get pg store client using transaction
	pgStoreClient, err := s.pgStore.NewClient(false)
	if err != nil {
		return err
	}

	// inserts task in pgstore
	err = pgStoreClient.UpdatePayment(ctx, reqPayment)
	if err != nil {
		return err
	}

	user, err := s.user.GetUserByID(ctx, reqPayment.UserID)
	if err != nil {
		return err
	}

	if reqPayment.Status == payment.StatusDone {
		user.Quota = 3
		user.Type = auth.TypePemium
	} else if reqPayment.Status == payment.StatusRejected {
		user.Quota = 0
		user.Type = auth.TypeFree
	}

	err = s.user.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
