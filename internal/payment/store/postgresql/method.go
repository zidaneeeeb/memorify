package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"hbdtoyou/internal/payment"

	"github.com/jmoiron/sqlx"
)

func (sc *storeClient) CreatePayment(ctx context.Context, reqPayment payment.Payment) (string, error) {
	// construct arguments filled with fields for the query
	argKV := map[string]interface{}{
		"user_id":           reqPayment.UserID,
		"content_id":        reqPayment.ContentID,
		"amount":            reqPayment.Amount,
		"proof_payment_url": reqPayment.ProofPaymentURL,
		"date":              reqPayment.Date,
		"status":            reqPayment.Status,
		"create_time":       reqPayment.CreateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryCreatePayment, argKV)
	if err != nil {
		return "", err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return "", err
	}
	query = sc.q.Rebind(query)

	// execute query
	var id string
	err = sc.q.QueryRowx(query, args...).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (sc *storeClient) GetPaymentByID(ctx context.Context, paymentID string) (payment.Payment, error) {
	query := fmt.Sprintf(queryGetPayment, "WHERE p.id = $1")

	// query single row
	var model paymentModel
	err := sc.q.QueryRowx(query, paymentID).StructScan(&model)
	if err != nil {
		if err == sql.ErrNoRows {
			return payment.Payment{}, payment.ErrDataNotFound
		}
		return payment.Payment{}, err
	}

	return model.format(), nil
}

func (sc *storeClient) GetPayments(ctx context.Context) ([]payment.Payment, error) {
	// define variables to custom query
	// argKV := make(map[string]interface{})
	// conditions := make([]string, 0)

	// if filter.UserID != "" {
	// 	id, err := uuid.Parse(filter.UserID)
	// 	if err != nil {
	// 		return nil, payment.ErrInvalidUserID
	// 	}
	// 	conditions = append(conditions, "p.user_id = :user_id")
	// 	argKV["user_id"] = id
	// }

	// if filter.Status > 0 {
	// 	conditions = append(conditions, "p.status = :status")
	// 	argKV["status"] = filter.Status
	// }

	// if filter.Status > 0 {
	// 	conditions = append(conditions, "p.status = :status")
	// 	argKV["status"] = filter.Status
	// }

	// construct strings to custom query
	// condition := strings.Join(conditions, " AND ")

	// // since the query does not contains "WHERE" yet, need
	// // to add it if needed
	// if len(conditions) > 0 {
	// 	condition = fmt.Sprintf("WHERE %s", condition)
	// }

	// construct query
	query := fmt.Sprintf(queryGetPayment, "ORDER BY p.date DESC")

	// prepare query
	query, args, err := sqlx.Named(query, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}
	query = sc.q.Rebind(query)

	// query to database
	rows, err := sc.q.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// read rows
	result := make([]payment.Payment, 0)
	for rows.Next() {
		var row paymentModel
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		result = append(result, row.format())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (sc *storeClient) UpdatePayment(ctx context.Context, reqPayment payment.Payment) error {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"id":                reqPayment.ID,
		"user_id":           reqPayment.UserID,
		"content_id":        reqPayment.ContentID,
		"amount":            reqPayment.Amount,
		"proof_payment_url": reqPayment.ProofPaymentURL,
		"date":              reqPayment.Date,
		"status":            reqPayment.Status,
		"update_time":       reqPayment.UpdateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryUpdatePayment, argsKV)
	if err != nil {
		return err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}
	query = sc.q.Rebind(query)

	// execute query
	_, err = sc.q.Exec(query, args...)
	return err
}
