package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

type Group struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

func (g *Group) String() string {
	out, err := json.Marshal(g)
	if err != nil {
		return fmt.Sprintf("%+v", g.Name)
	}
	return string(out)
}
