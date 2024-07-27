package repositories

import (
	"encoding/json"
	"fmt"
)

type Category struct {
	ID          uint
	Name        string
	Icon        string
	Color       string
	Description string
	Order       uint
}

func (c *Category) String() string {
	out, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("%+v", c.Name)
	}
	return string(out)
}
