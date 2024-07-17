package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/prulloac/fineasy/pkg"
)

type User struct {
	pkg.Model
	Hash              string
	Username          string
	Email             string
	ValidatedAt       sql.NullTime
	Disabled          bool
	InternalLoginData InternalLogin
	ExternalLoginData ExternalLogin
}

func (u *User) String() string {
	out, err := json.Marshal(u)
	if err != nil {
		return fmt.Sprintf("%+v", u.Username)
	}
	return string(out)
}

type InternalLogin struct {
	pkg.Model
	UserID                uint
	Password              string
	PasswordSalt          string
	Algorithm             pkg.Algorithm
	PasswordLastUpdatedAt time.Time
	LoginAttempts         int
	LastLoginAttempt      time.Time
	LastLoginSuccess      time.Time
}

func (i *InternalLogin) String() string {
	out, err := json.Marshal(i)
	if err != nil {
		return fmt.Sprintf("%+v", i.UserID)
	}
	return string(out)
}

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

type UserSession struct {
	pkg.Model
	UserID       uint
	LoginIP      string
	UserAgent    string
	LoggedInAt   time.Time
	LoggedOutAt  sql.NullTime
	SessionToken string
}

func (u *UserSession) String() string {
	out, err := json.Marshal(u)
	if err != nil {
		return fmt.Sprintf("%+v", u.ID)
	}
	return string(out)
}
