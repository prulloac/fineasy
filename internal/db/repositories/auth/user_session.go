package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/prulloac/fineasy/pkg"
)

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

func (a *AuthRepository) LogUserSession(uid uint, ip string, userAgent string) (*UserSession, error) {
	token, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	var session UserSession
	err = a.Persistence.QueryRow(`
	INSERT INTO user_sessions
	(user_id, login_ip, user_agent, session_token) VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, login_ip, user_agent, session_token, logged_in_at, logged_out_at, created_at, updated_at
	`, uid, ip, userAgent, token.String()).Scan(&session.ID, &session.UserID, &session.LoginIP, &session.UserAgent, &session.SessionToken, &session.LoggedInAt, &session.LoggedOutAt, &session.CreatedAt, &session.UpdatedAt)
	return &session, err
}
