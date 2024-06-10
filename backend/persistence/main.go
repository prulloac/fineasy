package persistence

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func Connection() *sql.DB {
	return connect()
}

func VerifySchema(db *sql.DB) {
	CreateCategoriesTable(db)
}
