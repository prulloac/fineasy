package repositories

import (
	"database/sql"
	"fmt"
	"os"
)

type TransactionsRepository struct {
	db *sql.DB
}

func NewTransactionsRepository(db *sql.DB) *TransactionsRepository {
	return &TransactionsRepository{db}
}

func (t *TransactionsRepository) CreateTable() {
	data, _ := os.ReadFile("internal/persistence/schema/transactions.sql")

	if data == nil {
		panic("Error reading transactions schema file!")
	}

	_, err := t.db.Exec(string(data))
	if err != nil {
		fmt.Println("Error creating transactions table!")
		panic(err)
	}
	fmt.Println("Transactions table created!")
}

func (t *TransactionsRepository) DropTable() {
	_, err := t.db.Exec("DROP TABLE IF EXISTS transactions")
	if err != nil {
		fmt.Println("Error dropping transactions table!")
		panic(err)
	}
	fmt.Println("Transactions table dropped!")
}
