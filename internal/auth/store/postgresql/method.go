package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"hbdtoyou/internal/auth"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func (sc *storeClient) CreateUser(ctx context.Context, reqUser auth.User) (string, error) {
	argKV := map[string]interface{}{
		"fullname":    reqUser.Fullname,
		"username":    reqUser.Username,
		"email":       reqUser.Email,
		"type":        reqUser.Type,
		"quota":       reqUser.Quota,
		"create_time": reqUser.CreateTime,
	}

	query, args, err := sqlx.Named(queryCreateUser, argKV)
	if err != nil {
		return "", err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return "", err
	}

	query = sc.q.Rebind(query)

	var userID string
	err = sc.q.QueryRowx(query, args...).Scan(&userID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr != nil {
			if pqErr.Code.Name() == "unique_violation" {
				return "", auth.ErrUserAlreadyExist
			}
		}
		return "", err
	}

	return userID, nil
}

func (s *storeClient) GetUserAuth(ctx context.Context, filter auth.GetUserAuthFilter) (auth.User, error) {
	argKV := make(map[string]interface{})
	conditions := make([]string, 0)

	conditions = append(conditions, "u.status != :status")
	argKV["status"] = "3"

	if filter.Email != "" {
		conditions = append(conditions, "u.email = :email")
		argKV["email"] = filter.Email
	}

	if filter.UserID != "" {
		conditions = append(conditions, "u.id = :id")
		argKV["id"] = filter.UserID
	}

	condition := strings.Join(conditions, " AND ")
	if len(conditions) > 0 {
		condition = fmt.Sprintf("WHERE %s", condition)
	}

	query := fmt.Sprintf(queryGetUser, condition)
	query, args, err := sqlx.Named(query, argKV)
	if err != nil {
		return auth.User{}, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return auth.User{}, err
	}

	query = s.q.Rebind(query)

	var userAuth UserModel
	if err := s.q.QueryRowx(query, args...).StructScan(&userAuth); err != nil {
		if err == sql.ErrNoRows {
			return auth.User{}, auth.ErrDataNotFound
		}

		return auth.User{}, err
	}

	return userAuth.format(), nil
}

func (sc *storeClient) UpdateUser(ctx context.Context, reqUser auth.User) error {
	argsKV := map[string]interface{}{
		"id":              reqUser.ID,
		"fullname":        reqUser.Fullname,
		"username":        reqUser.Username,
		"email":           reqUser.Email,
		"type":            reqUser.Type,
		"quota":           reqUser.Quota,
		"update_time":     reqUser.UpdateTime,
	}

	query, args, err := sqlx.Named(queryUpdateUser, argsKV)
	if err != nil {
		return err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}

	query = sc.q.Rebind(query)

	_, err = sc.q.Exec(query, args...)
	return err
}