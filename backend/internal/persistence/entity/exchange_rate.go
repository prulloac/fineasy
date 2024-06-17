package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

type ExchangeRate struct {
	ID         int       `json:"id"`
	CurrencyID int       `json:"currency_id"`
	GroupID    int       `json:"group_id"`
	Rate       float64   `json:"rate"`
	Date       time.Time `json:"date"`
}

func (e *ExchangeRate) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("%+v", e.ID)
	}
	return string(out)
}
