package preferences

import (
	"encoding/json"
	"time"
)

type UserData struct {
	ID          uint
	UserID      uint
	AvatarURL   string
	DisplayName string
	Currency    string
	Language    string
	Timezone    string
	UpsertedAt  time.Time
}

func (u *UserData) String() string {
	out, err := json.Marshal(u)
	if err != nil {
		return u.DisplayName
	}
	return string(out)
}

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
