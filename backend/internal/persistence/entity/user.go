package entity

import (
	"encoding/json"
	"fmt"
	"time"
)

type User struct {
	ID        int
	Username  string
	Email     string
	CreatedAt time.Time
	UpdateAt  time.Time
}

func (u *User) String() string {
	out, err := json.Marshal(u)
	if err != nil {
		return fmt.Sprintf("%+v", u.Username)
	}
	return string(out)
}
