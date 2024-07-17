package transactions

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/prulloac/fineasy/pkg"
)

type Account struct {
	ID        uint
	CreatedBy uint
	GroupID   uint
	Currency  string
	Balance   float64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *Account) String() string {
	out, err := json.Marshal(a)
	if err != nil {
		return fmt.Sprintf("%+v", a.Name)
	}
	return string(out)
}

type Budget struct {
	ID        uint
	Name      string
	AccountID uint
	Currency  string
	Amount    float64
	CreatedBy uint
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *Budget) String() string {
	out, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("%+v", b.Name)
	}
	return string(out)
}

type Transaction struct {
	ID              uint
	Category        string
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
