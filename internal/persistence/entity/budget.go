package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

type Budget struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	AccountID  int       `json:"account_id"`
	CurrencyID int       `json:"currency_id"`
	Amount     float32   `json:"amount"`
	CreatedBy  int       `json:"created_by"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	CreatedAt  time.Time `json:"created_at"`
	UpdateAt   time.Time `json:"updated_at"`
}

func (b *Budget) String() string {
	out, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("%+v", b.Name)
	}
	return string(out)
}
