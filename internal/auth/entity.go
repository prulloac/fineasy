package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type User struct {
	ID                int           `json:"id" validate:"required,min=1"`
	Hash              string        `json:"hash" validate:"required,uuid7"`
	Username          string        `json:"username" validate:"required,min=1"`
	Email             string        `json:"email" validate:"required,email"`
	ValidatedAt       sql.NullTime  `json:"validated_at"`
	Disabled          bool          `json:"disabled"`
	CreatedAt         time.Time     `json:"created_at" validate:"required,past_time"`
	UpdatedAt         time.Time     `json:"updated_at" validate:"required,past_time"`
	internalLoginData InternalLogin `json:"-"`
	externalLoginData ExternalLogin `json:"-"`
}

func (u *User) String() string {
	out, err := json.Marshal(u)
	if err != nil {
		return fmt.Sprintf("%+v", u.Username)
	}
	return string(out)
}

type Algorithm uint16

const (
	None Algorithm = iota
	SHA256
	SHA512
	SHA3_256
	SHA3_512
	Base64
)

func (a Algorithm) Name() string {
	switch a {
	case SHA256:
		return "SHA256"
	case SHA512:
		return "SHA512"
	case SHA3_256:
		return "SHA3_256"
	case SHA3_512:
		return "SHA3_512"
	case Base64:
		return "Base64"
	case None:
		return ""
	default:
		panic("invalid algorithm")
	}
}

type TokenTypes uint16

const (
	AccessToken TokenTypes = iota
	RefreshToken
	PasswordResetToken
	EmailVerificationToken
)

func (t TokenTypes) Name() string {
	switch t {
	case AccessToken:
		return "AccessToken"
	case RefreshToken:
		return "RefreshToken"
	case PasswordResetToken:
		return "PasswordResetToken"
	case EmailVerificationToken:
		return "EmailVerificationToken"
	default:
		panic("invalid token type")
	}
}

type InternalLogin struct {
	ID                    int       `json:"id" validate:"required,min=1"`
	UserID                int       `json:"user_id" validate:"required,min=1"`
	Email                 string    `json:"email" validate:"required,email"`
	Password              string    `json:"password" validate:"required,min=1"`
	PasswordSalt          string    `json:"password_salt" validate:"required"`
	Algorithm             Algorithm `json:"algorithm" validate:"required,min=0"`
	PasswordLastUpdatedAt time.Time `json:"password_last_updated_at" validate:"required,past_time"`
	LoginAttempts         int       `json:"login_attempts" validate:"required,min=1"`
	LastLoginAttempt      time.Time `json:"last_login_attempt" validate:"required,past_time"`
	LastLoginSuccess      time.Time `json:"last_login_success" validate:"required,past_time"`
	CreatedAt             time.Time `json:"created_at" validate:"required,past_time"`
	UpdatedAt             time.Time `json:"updated_at" validate:"required,past_time"`
}

func (i *InternalLogin) String() string {
	out, err := json.Marshal(i)
	if err != nil {
		return fmt.Sprintf("%+v", i.Email)
	}
	return string(out)
}

type LoginToken struct {
	ID        int       `json:"id" validate:"required,min=1"`
	UserID    int       `json:"user_id" validate:"required,min=1"`
	Token     string    `json:"token" validate:"required,min=1"`
	TokenType int       `json:"token_type" validate:"required,min=1"`
	ExpiresAt time.Time `json:"expires_at" validate:"required,past_time"`
	UsedAt    time.Time `json:"used_at" validate:"required,past_time"`
	CreatedAt time.Time `json:"created_at" validate:"required,past_time"`
}

func (l *LoginToken) String() string {
	out, err := json.Marshal(l)
	if err != nil {
		return fmt.Sprintf("%+v", l.Token)
	}
	return string(out)
}

type ExternalLoginProvider struct {
	ID        int       `json:"id" validate:"required,min=1"`
	Name      string    `json:"name" validate:"required,min=2,max=254"`
	Type      int       `json:"type" validate:"required"`
	Endpoint  string    `json:"endpoint" validate:"required,url"`
	Enabled   bool      `json:"enabled" validate:"required,boolean"`
	CreatedAt time.Time `json:"created_at" validate:"required,past_time"`
	UpdatedAt time.Time `json:"updated_at" validate:"required,past_time"`
}

func (e *ExternalLoginProvider) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("%+v", e.Name)
	}
	return string(out)
}

type ExternalLogin struct {
	ID         int       `json:"id" validate:"required,min=1"`
	UserID     int       `json:"user_id" validate:"required,min=1"`
	ProviderID int       `json:"provider_id" validate:"required,min=1"`
	CreatedAt  time.Time `json:"created_at" validate:"required,past_time"`
	UpdatedAt  time.Time `json:"updated_at" validate:"required,past_time"`
}

func (e *ExternalLogin) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("%+v", e.ProviderID)
	}
	return string(out)
}

type ExternalLoginToken struct {
	ID         int       `json:"id" validate:"required,min=1"`
	UserID     int       `json:"user_id" validate:"required,min=1"`
	ProviderID int       `json:"provider_id" validate:"required,min=1"`
	LoginIP    string    `json:"login_ip" validate:"required,min=1"`
	UserAgent  string    `json:"user_agent" validate:"required,min=1"`
	LoggedInAt time.Time `json:"logged_in_at" validate:"required,past_time"`
	Token      string    `json:"token" validate:"required,min=1"`
	CreatedAt  time.Time `json:"created_at" validate:"required,past_time"`
}

func (e *ExternalLoginToken) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("%+v", e.Token)
	}
	return string(out)
}

type UserSession struct {
	ID          int       `json:"id" validate:"required,min=1"`
	UserID      int       `json:"user_id" validate:"required,min=1"`
	LoginIP     string    `json:"login_ip" validate:"required,min=1"`
	UserAgent   string    `json:"user_agent" validate:"required,min=1"`
	LoggedInAt  time.Time `json:"logged_in_at" validate:"required,past_time"`
	LoggedOutAt time.Time `json:"logged_out_at" validate:"required,past_time"`
	CreatedAt   time.Time `json:"created_at" validate:"required,past_time"`
	UpdatedAt   time.Time `json:"updated_at" validate:"required,past_time"`
}

func (u *UserSession) String() string {
	out, err := json.Marshal(u)
	if err != nil {
		return fmt.Sprintf("%+v", u.ID)
	}
	return string(out)
}
