package entity

import (
	"encoding/json"
	"fmt"
)

type UserGroup struct {
	ID      int `json:"id"`
	UserID  int `json:"user_id"`
	GroupID int `json:"group_id"`
}

func (ug *UserGroup) String() string {
	out, err := json.Marshal(ug)
	if err != nil {
		return fmt.Sprintf("%+v", ug.UserID)
	}
	return string(out)
}
