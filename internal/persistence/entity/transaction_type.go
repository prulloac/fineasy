package entity

import (
	"encoding/json"
	"fmt"
)

type TransactionType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (t *TransactionType) String() string {
	out, err := json.Marshal(t)
	if err != nil {
		return fmt.Sprintf("%+v", t.Name)
	}
	return string(out)
}

func (t *TransactionType) IsIncome() bool {
	return t.ID == 1
}

func (t *TransactionType) IsExpense() bool {
	return t.ID == 2
}

func (t *TransactionType) IsSaving() bool {
	return t.ID == 3
}

func Income() TransactionType {
	return TransactionType{ID: 1, Name: "Income"}
}

func Expense() TransactionType {
	return TransactionType{ID: 2, Name: "Expense"}
}

func Saving() TransactionType {
	return TransactionType{ID: 3, Name: "Saving"}
}

func TransactionTypeByID(id int) (TransactionType, error) {
	switch id {
	case 1:
		return Income(), nil
	case 2:
		return Expense(), nil
	case 3:
		return Saving(), nil
	default:
		return TransactionType{}, fmt.Errorf("unknown transaction type id: %d", id)
	}
}
