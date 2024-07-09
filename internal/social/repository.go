package social

import "database/sql"

type SocialRepository struct {
	DB *sql.DB
}

func NewSocialRepository(db *sql.DB) *SocialRepository {
	return &SocialRepository{DB: db}
}

func (s *SocialRepository) CreateTable() error {
	_, err := s.DB.Exec(`
	CREATE TABLE IF NOT EXISTS friends (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		friend_id INT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		relation_type INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS friend_requests (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		friend_id INT NOT NULL,
		status INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS groups (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		created_by INT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS user_groups (
		id SERIAL PRIMARY KEY,
		group_id INT NOT NULL,
		user_id INT NOT NULL,
		joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		left_at TIMESTAMP
	);
	`)
	return err
}

func (s *SocialRepository) DropTable() error {
	_, err := s.DB.Exec(`
	DROP TABLE IF EXISTS user_groups;
	DROP TABLE IF EXISTS groups;
	DROP TABLE IF EXISTS friend_requests;
	DROP TABLE IF EXISTS friends;
	`)
	return err
}

func (s *SocialRepository) AddFriend(userID, friendID int) (*Friend, error) {
	_, err := s.DB.Exec(`
	INSERT INTO friends (user_id, friend_id, relation_type)
	VALUES ($1, $2, $3)
	`, userID, friendID, 1)
	if err != nil {
		return nil, err
	}
	return &Friend{
		UserID:       userID,
		FriendID:     friendID,
		RelationType: Contact,
	}, nil
}
