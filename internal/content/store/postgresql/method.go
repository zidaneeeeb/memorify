package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"hbdtoyou/internal/content"
	"strings"

	contextlib "hbdtoyou/pkg/context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (sc *storeClient) CreateContent(ctx context.Context, reqContent content.Content) (string, error) {
	// construct arguments filled with fields for the query
	argKV := map[string]interface{}{
		"user_id":                  reqContent.UserID,
		"template_id":              reqContent.TemplateID,
		"detail_content_json_text": reqContent.DetailContentJSONText,
		"status":                   reqContent.Status,
		"create_time":              reqContent.CreateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryCreateContent, argKV)
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

func (sc *storeClient) GetContents(ctx context.Context, filter content.GetContentsFilter) ([]content.Content, error) {
	// define variables to custom query
	argKV := make(map[string]interface{})
	conditions := make([]string, 0)

	if filter.UserID != "" {
		id, err := uuid.Parse(filter.UserID)
		if err != nil {
			return nil, content.ErrInvalidUserID
		}
		conditions = append(conditions, "c.user_id = :user_id")
		argKV["user_id"] = id
	}

	if filter.TemplateID != "" {
		conditions = append(conditions, "c.template_id = :template_id")
		argKV["template_id"] = filter.TemplateID
	}

	if filter.TemplateLabel != "" {
		conditions = append(conditions, "t.label = :template_label")
		argKV["template_label"] = filter.TemplateLabel
	}

	// if filter.Status > 0 {
	// 	conditions = append(conditions, "c.status = :status")
	// 	argKV["status"] = filter.Status
	// }

	// construct strings to custom query
	condition := strings.Join(conditions, " AND ")

	// since the query does not contains "WHERE" yet, need
	// to add it if needed
	if len(conditions) > 0 {
		condition = fmt.Sprintf("WHERE %s", condition)
	}

	// construct query
	query := fmt.Sprintf(queryGetContent, condition)

	// prepare query
	query, args, err := sqlx.Named(query, argKV)
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
	result := make([]content.Content, 0)
	for rows.Next() {
		var row contentModel
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

func (sc *storeClient) GetContentByID(ctx context.Context, contentID string) (content.Content, error) {
	query := fmt.Sprintf(queryGetContent, "WHERE c.id = $1")

	// query single row
	var model contentModel
	err := sc.q.QueryRowx(query, contentID).StructScan(&model)
	if err != nil {
		if err == sql.ErrNoRows {
			return content.Content{}, content.ErrDataNotFound
		}
		return content.Content{}, err
	}

	return model.format(), nil
}

func (sc *storeClient) UpdateContent(ctx context.Context, reqContent content.Content) error {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"id":                       reqContent.ID,
		"template_id":              reqContent.TemplateID,
		"detail_content_json_text": reqContent.DetailContentJSONText,
		"status":                   reqContent.Status,
		"update_time":              reqContent.UpdateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryUpdateContent, argsKV)
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

func (sc *storeClient) DeleteContentByID(ctx context.Context, contentID string) error {
	// get user ID
	userID, ok := contextlib.GetUserID(ctx)
	if !ok {
		return content.ErrInvalidUserID
	}

	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"id":      contentID,
		"user_id": userID,
	}

	// prepare query
	query, args, err := sqlx.Named(queryDeleteContent, argsKV)
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
	if err != nil {
		return err
	}

	return nil
}
