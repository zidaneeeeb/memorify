package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"hbdtoyou/internal/template"
	"strings"

	"github.com/jmoiron/sqlx"
)

func (sc *storeClient) CreateTemplate(ctx context.Context, reqTemplate template.Template) (string, error) {
	// construct arguments filled with fields for the query
	argKV := map[string]interface{}{
		"name":          reqTemplate.Name,
		"label":         reqTemplate.Label,
		"thumbnail_uri": reqTemplate.ThumbnailURI,
		"create_time":   reqTemplate.CreateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryCreateTemplate, argKV)
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

func (sc *storeClient) GetTemplateByID(ctx context.Context, templateID string) (template.Template, error) {
	query := fmt.Sprintf(queryGetTemplate, "WHERE t.id = $1")

	// query single row
	var model templateModel
	err := sc.q.QueryRowx(query, templateID).StructScan(&model)
	if err != nil {
		if err == sql.ErrNoRows {
			return template.Template{}, template.ErrTemplateNotFound
		}
		return template.Template{}, err
	}

	return model.format(), nil
}

func (sc *storeClient) GetTemplates(ctx context.Context, filter template.GetTemplatesFilter) ([]template.Template, error) {
	// define variables to custom query
	argKV := make(map[string]interface{})
	conditions := make([]string, 0)

	if filter.Label > 0 {
		conditions = append(conditions, "t.label = :label")
		argKV["label"] = filter.Label
	}

	// construct strings to custom query
	condition := strings.Join(conditions, " AND ")

	// since the query does not contains "WHERE" yet, need
	// to add it if needed
	if len(conditions) > 0 {
		condition = fmt.Sprintf("WHERE %s", condition)
	}

	// construct query
	query := fmt.Sprintf(queryGetTemplate, condition)

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
	result := make([]template.Template, 0)
	for rows.Next() {
		var row templateModel
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

func (sc *storeClient) UpdateTemplate(ctx context.Context, reqTemplate template.Template) error {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"id":            reqTemplate.ID,
		"name":          reqTemplate.Name,
		"label":         reqTemplate.Label,
		"thumbnail_uri": reqTemplate.ThumbnailURI,
		"update_time":   reqTemplate.UpdateTime,
	}

	// prepare query
	query, args, err := sqlx.Named(queryUpdateTemplate, argsKV)
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

func (sc *storeClient) DeleteTemplateByID(ctx context.Context, templateID string) error {
	// construct arguments filled with fields for the query
	argsKV := map[string]interface{}{
		"id": templateID,
	}

	// prepare query
	query, args, err := sqlx.Named(queryDeleteTemplate, argsKV)
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
