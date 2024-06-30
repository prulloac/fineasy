package transactions

import (
	"database/sql"
)

type TransactionsRepository struct {
	db *sql.DB
}

func NewTransactionsRepository(db *sql.DB) *TransactionsRepository {
	return &TransactionsRepository{db: db}
}

func (r *TransactionsRepository) CreateTable() error {
	return nil
}

func (r *TransactionsRepository) DropTable() error {
	return nil
}
