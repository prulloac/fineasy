package preferences

import (
	p "github.com/prulloac/fineasy/internal/persistence"
)

type Repository struct {
	Persistence *p.Persistence
}

func NewRepository(persistence *p.Persistence) *Repository {
	return &Repository{persistence}
}

func (r *Repository) CreateTables() error {
	_, err := r.Persistence.Exec(`
		CREATE TABLE IF NOT EXISTS user_preferences (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			key VARCHAR(255) NOT NULL,
			value TEXT NOT NULL,
			upserted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE UNIQUE INDEX IF NOT EXISTS idx_user_preferences_user_id_key ON user_preferences (user_id, key);

		CREATE TABLE IF NOT EXISTS user_data (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			avatar_url TEXT NOT NULL,
			display_name VARCHAR(255) NOT NULL,
			currency VARCHAR(6) NOT NULL,
			language VARCHAR(6) NOT NULL,
			timezone VARCHAR(255) NOT NULL,
			upserted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE UNIQUE INDEX IF NOT EXISTS idx_user_data_user_id ON user_data (user_id);
	`)
	return err
}

func (r *Repository) Close() {
	r.Persistence.Close()
}

func (r *Repository) DropTables() error {
	_, err := r.Persistence.Exec(`
		DROP TABLE IF EXISTS user_preferences;
		DROP TABLE IF EXISTS user_data;
	`)
	return err
}

func (r *Repository) CreateUserData(userID uint) (*UserData, error) {
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

func (r *Repository) GetUserDataByUserID(userID uint) (*UserData, error) {
	row := r.Persistence.QueryRow(`
		SELECT id, user_id, avatar_url, display_name, currency, language, timezone
		FROM user_data
		WHERE user_id = $1
	`, userID)
	userData := UserData{}
	err := row.Scan(&userData.ID, &userData.UserID, &userData.AvatarURL, &userData.DisplayName, &userData.Currency, &userData.Language, &userData.Timezone)
	if err != nil {
		return nil, err
	}
	return &userData, nil
}
