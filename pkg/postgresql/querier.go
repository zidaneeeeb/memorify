package postgresql

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -destination ./mock_querier.go -package postgresql Querier

// Querier provides mechanism to do query.
//
// Querier is the abstraction of sqlx.DB and sqlx.TX object.
// Querier can be added to an object to supports both single
// and transaction query execution.
type Querier interface {
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Rebind(query string) string
	Exec(query string, args ...interface{}) (sql.Result, error)
}
