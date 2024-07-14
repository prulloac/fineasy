package transactions

import (
	"encoding/json"
	"fmt"
	"time"
)

type Group struct {
	ID        uint      `json:"id" validate:"required,min=1"`
	Name      string    `json:"name" validate:"required,min=1"`
	CreatedBy uint      `json:"created_by" validate:"required,min=1"`
	CreatedAt time.Time `json:"created_at" validate:"required,past_time"`
	UpdatedAt time.Time `json:"updated_at" validate:"required,past_time"`
	UserCount uint      `json:"user_count" validate:"required,min=0"`
}

func (g *Group) String() string {
	out, err := json.Marshal(g)
	if err != nil {
		return fmt.Sprintf("%+v", g.Name)
	}
	return string(out)
}

type UserGroup struct {
	ID        uint      `json:"id" validate:"required,min=1"`
	UserID    uint      `json:"user_id" validate:"required,min=1"`
	GroupID   uint      `json:"group_id" validate:"required,min=1"`
	CreatedAt time.Time `json:"created_at" validate:"required,past_time"`
}

func (u *UserGroup) String() string {
	out, err := json.Marshal(u)
	if err != nil {
		return fmt.Sprintf("%+v", u.GroupID)
	}
	return string(out)
}

type Account struct {
	ID        uint      `json:"id" validate:"required,min=1"`
	CreatedBy uint      `json:"created_by" validate:"required,min=1"`
	GroupID   uint      `json:"group_id" validate:"required,min=1"`
	Currency  string    `json:"currency" validate:"required,min=1"`
	Balance   float64   `json:"balance" validate:"required,min=0"`
	Name      string    `json:"name" validate:"required,min=1"`
	CreatedAt time.Time `json:"created_at" validate:"required,past_time"`
	UpdatedAt time.Time `json:"updated_at" validate:"required,past_time"`
}

func (a *Account) String() string {
	out, err := json.Marshal(a)
	if err != nil {
		return fmt.Sprintf("%+v", a.Name)
	}
	return string(out)
}

type Budget struct {
	ID        uint      `json:"id" validate:"required,min=1"`
	Name      string    `json:"name" validate:"required,min=1"`
	AccountID uint      `json:"account_id" validate:"required,min=1"`
	Currency  string    `json:"currency" validate:"required,min=1"`
	Amount    float64   `json:"amount" validate:"required,min=0"`
	CreatedBy uint      `json:"created_by" validate:"required,min=1"`
	StartDate time.Time `json:"start_date" validate:"required,past_time"`
	EndDate   time.Time `json:"end_date" validate:"required,past_time"`
	CreatedAt time.Time `json:"created_at" validate:"required,past_time"`
	UpdatedAt time.Time `json:"updated_at" validate:"required,past_time"`
}

func (b *Budget) String() string {
	out, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("%+v", b.Name)
	}
	return string(out)
}

type Category struct {
	ID          uint   `json:"id" validate:"required,min=1"`
	Name        string `json:"name" validate:"required,min=1"`
	Icon        string `json:"icon" validate:"required,min=1"`
	Color       string `json:"color" validate:"required,min=1"`
	Description string `json:"description" validate:"required,min=1"`
	Order       uint   `json:"order" validate:"required,min=0"`
}

func (c *Category) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("%+v", c.Name)
	}
	return string(out)
}

type transactionType uint8

type Transaction struct {
	ID              uint            `json:"id" validate:"required,min=1"`
	CategoryID      uint            `json:"category_id" validate:"required,min=1"`
	Currency        string          `json:"currency" validate:"required,min=1"`
	CurrencyRate    float64         `json:"currency_rate" validate:"required,min=0"`
	TransactionType transactionType `json:"transaction_type" validate:"required,min=1"`
	BudgetID        uint            `json:"budget_id" validate:"required,min=1"`
	Amount          float64         `json:"amount" validate:"required,min=0"`
	Date            time.Time       `json:"date" validate:"required,past_time"`
	ExecutedByName  int             `json:"executed_by" validate:"required,min=1"`
	ExecutedByID    uint            `json:"executed_by_id"`
	Description     string          `json:"description" validate:"required,min=1"`
	ReceiptURL      string          `json:"receipt_url" validate:"required,min=1"`
	RegisteredBy    uint            `json:"registered_by" validate:"required,min=1"`
	RegisteredAt    time.Time       `json:"registered_at" validate:"required,past_time"`
	CreatedAt       time.Time       `json:"created_at" validate:"required,past_time"`
	UpdatedAt       time.Time       `json:"updated_at" validate:"required,past_time"`
}

func (t *Transaction) String() string {
	out, err := json.Marshal(t)
	if err != nil {
		return fmt.Sprintf("%+v", t.Description)
	}
	return string(out)
}
