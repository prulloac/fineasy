package repositories

import (
	"encoding/json"
	"fmt"
	"time"
)

type ExternalLogin struct {
	ID         int
	UserID     int
	ProviderID int
	CreatedAt  time.Time
}

func (e *ExternalLogin) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("%+v", e.ProviderID)
	}
	return string(out)
}
