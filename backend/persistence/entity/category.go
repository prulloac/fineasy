package entity

import (
	"encoding/json"
	"fmt"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Order       int    `json:"order"`
	GroupID     int    `json:"group_id"`
}

func (c *Category) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("%+v", c.Name)
	}
	return string(out)
}
