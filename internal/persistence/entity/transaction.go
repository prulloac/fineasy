package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

type Transaction struct {
	ID              int       `json:"id"`
	CategoryID      int       `json:"category_id"`
	CurrencyID      int       `json:"currency_id"`
	TransactionType int       `json:"transaction_type"`
	AccountID       int       `json:"account_id"`
	Amount          float32   `json:"amount"`
	Date            time.Time `json:"date"`
	ExecutedBy      int       `json:"executed_by"`
	Description     string    `json:"description"`
	ReceiptURL      string    `json:"receipt_url"`
	RegisteredAt    time.Time `json:"registered_at"`
	RegisteredBy    int       `json:"registered_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (t *Transaction) String() string {
	out, err := json.Marshal(t)
	if err != nil {
		return fmt.Sprintf("%+v", t.ID)
	}
	return string(out)
}
