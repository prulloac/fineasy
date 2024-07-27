package repositories

import (
	"encoding/json"
	"fmt"
	"time"
)

type ExternalLoginToken struct {
	ID              int
	ExternalLoginID int
	LoginIP         string
	UserAgent       string
	LoggedInAt      time.Time
	Token           string
	CreatedAt       time.Time
}

func (e *ExternalLoginToken) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("%+v", e.Token)
	}
	return string(out)
}
