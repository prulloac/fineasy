package repositories

import (
	"encoding/json"
	"time"
)

type UserPreferences struct {
	ID         uint
	UserID     uint
	Key        string
	Value      string
	UpsertedAt time.Time
}

func (u *UserPreferences) String() string {
	out, err := json.Marshal(u)
	if err != nil {
		return u.Key
	}
	return string(out)
}
