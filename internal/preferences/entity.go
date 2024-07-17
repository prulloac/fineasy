package userpreferences

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/prulloac/fineasy/pkg"
)

type Entry struct {
	Key   string
	Value interface{}
}

type UserPreference struct {
	ID         int       `json:"id" validate:"required,min=1"`
	UserID     int       `json:"user_id" validate:"required,min=1"`
	Theme      string    `json:"theme" validate:"required,min=1"`
	Language   string    `json:"language" validate:"required,min=1"`
	Entries    []Entry   `json:"entries" validate:"required,min=1"`
	UpsertedAt time.Time `json:"upserted_at" validate:"required,past_time"`
}

func (u *UserPreference) String() string {
	out, err := json.Marshal(u)
	if err != nil {
		return fmt.Sprintf("%+v", u.Theme)
	}
	return string(out)
}

type GroupPreference struct {
	ID         int       `json:"id" validate:"required,min=1"`
	GroupID    int       `json:"group_id" validate:"required,min=1"`
	Theme      string    `json:"theme" validate:"required,min=1"`
	Language   string    `json:"language" validate:"required,min=1"`
	Entries    []Entry   `json:"entries" validate:"required,min=1"`
	UpsertedAt time.Time `json:"upserted_at" validate:"required,past_time"`
}

func (g *GroupPreference) String() string {
	out, err := json.Marshal(g)
	if err != nil {
		return fmt.Sprintf("%+v", g.Theme)
	}
	return string(out)
}

type Category struct {
	pkg.Model
	UserID      uint   `json:"user_id" validate:"min=1"`
	GroupID     uint   `json:"group_id" validate:"min=1"`
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
