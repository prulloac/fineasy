package repositories

import (
	"encoding/json"
	"time"
)

type UserData struct {
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

func (r *CoreRepository) CreateUserData(userID uint) (*UserData, error) {
	_, err := r.Persistence.Exec(`
		INSERT INTO user_data (user_id, avatar_url, display_name, currency, language, timezone)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO NOTHING
	`, userID, "", "", "USD", "en", "UTC")
	if err != nil {
		return nil, err
	}
	return r.GetUserDataByUserID(userID)
}

func (r *CoreRepository) GetUserDataByUserID(userID uint) (*UserData, error) {
	row := r.Persistence.QueryRow(`
		SELECT user_id, avatar_url, display_name, currency, language, timezone
		FROM user_data
		WHERE user_id = $1
	`, userID)
	userData := UserData{}
	err := row.Scan(&userData.UserID, &userData.AvatarURL, &userData.DisplayName, &userData.Currency, &userData.Language, &userData.Timezone)
	if err != nil {
		return nil, err
	}
	return &userData, nil
}
