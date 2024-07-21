package repositories

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/prulloac/fineasy/pkg"
)

type LoginToken struct {
	ID        int
	UserID    int
	Token     string
	TokenType pkg.TokenType
	ExpiresAt time.Time
	UsedAt    time.Time
	CreatedAt time.Time
}

func (l *LoginToken) String() string {
	out, err := json.Marshal(l)
	if err != nil {
		return fmt.Sprintf("%+v", l.Token)
	}
	return string(out)
}
