package repositories

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/prulloac/fineasy/pkg"
)

type Transaction struct {
	ID              uint
	CategoryID      string
	Currency        string
	CurrencyRate    float64
	TransactionType pkg.TransactionType
	BudgetID        uint
	Amount          float64
	Date            time.Time
	ExecutedByName  int
	ExecutedByID    uint
	Description     string
	ReceiptURL      string
	RegisteredBy    uint
	RegisteredAt    time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (t *Transaction) String() string {
	out, err := json.Marshal(t)
	if err != nil {
		return fmt.Sprintf("%+v", t.Description)
	}
	return string(out)
}
