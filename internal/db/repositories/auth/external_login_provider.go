package repositories

import (
	"encoding/json"
	"fmt"
	"time"
)

type ExternalLoginProvider struct {
	ID        int
	Name      string
	Type      int
	Endpoint  string
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (e *ExternalLoginProvider) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("%+v", e.Name)
	}
	return string(out)
}
