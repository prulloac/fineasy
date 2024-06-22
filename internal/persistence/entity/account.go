package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

type Account struct {
	ID         int       `json:"id"`
	CreatedBy  int       `json:"created_by"`
	GroupID    int       `json:"group_id"`
	CurrencyID int       `json:"currency_id"`
	Balance    float32   `json:"balance"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Disabled   bool      `json:"disabled"`
}

func (c *Account) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("%+v", c.Name)
	}
	return string(out)
}
