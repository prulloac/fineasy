package social

import (
	"database/sql"

	"github.com/prulloac/fineasy/pkg"
)

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

func (s *SocialRepository) AddFriend(userID, friendID int) (*FriendRequestOutput, error) {
	var f FriendRequest
	err := s.DB.QueryRow(`
	INSERT INTO friend_requests (user_id, friend_id, status)
	VALUES ($1, $2, $3)
	RETURNING id, user_id, friend_id, status, created_at, updated_at
	`, userID, friendID, Pending).Scan(&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &FriendRequestOutput{
		UserID:   f.UserID,
		FriendID: f.FriendID,
		Status:   f.Status.String(),
	}, nil
}

func (s *SocialRepository) GetFriends(userID int) ([]FriendShipOutput, error) {
	rows, err := s.DB.Query(`
	SELECT user_id, friend_id, relation_type
	FROM friends
	WHERE user_id = $1 OR friend_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	friends := []FriendShipOutput{}
	for rows.Next() {
		var f Friend
		if err := rows.Scan(&f.UserID, &f.FriendID, &f.RelationType); err != nil {
			return nil, err
		}
		if err = pkg.ValidateStruct(f); err != nil {
			return nil, err
		}
		friends = append(friends, FriendShipOutput{
			UserID:       f.UserID,
			FriendID:     f.FriendID,
			RelationType: f.RelationType.String(),
		})
	}
	return friends, nil
}

func (s *SocialRepository) GetFriendRequests(userID int) ([]FriendRequestOutput, error) {
	rows, err := s.DB.Query(`
	SELECT user_id, friend_id, status
	FROM friend_requests
	WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []FriendRequestOutput{}
	for rows.Next() {
		var f FriendRequest
		if err := rows.Scan(&f.UserID, &f.FriendID, &f.Status); err != nil {
			return nil, err
		}
		if err = pkg.ValidateStruct(f); err != nil {
			return nil, err
		}
		requests = append(requests, FriendRequestOutput{
			UserID:   f.UserID,
			FriendID: f.FriendID,
			Status:   f.Status.String(),
		})
	}
	return requests, nil
}

func (s *SocialRepository) AcceptFriendRequest(userID, friendID int) (*FriendRequestOutput, error) {
	var f FriendRequest
	err := s.DB.QueryRow(`
	UPDATE friend_requests
	SET status = $3
	WHERE user_id = $1 AND friend_id = $2
	RETURNING id, user_id, friend_id, status, created_at, updated_at
	`, userID, friendID, Accepted).Scan(&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	// Add friend to friends table
	_, err = s.DB.Exec(`
	INSERT INTO friends (user_id, friend_id, relation_type)
	VALUES ($1, $2, $3)
	`, userID, friendID, Contact)
	if err != nil {
		return nil, err
	}
	return &FriendRequestOutput{
		UserID:   f.UserID,
		FriendID: f.FriendID,
		Status:   f.Status.String(),
	}, nil
}

func (s *SocialRepository) RejectFriendRequest(userID, friendID int) (*FriendRequestOutput, error) {
	var f FriendRequest
	err := s.DB.QueryRow(`
	UPDATE friend_requests
	SET status = $3
	WHERE user_id = $1 AND friend_id = $2
	RETURNING id, user_id, friend_id, status, created_at, updated_at
	`, userID, friendID, Declined).Scan(&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &FriendRequestOutput{
		UserID:   f.UserID,
		FriendID: f.FriendID,
		Status:   f.Status.String(),
	}, nil
}

func (s *SocialRepository) CreateGroup(name string, members []int, createdBy int) (*GroupOutput, error) {
	var g Group
	err := s.DB.QueryRow(`
	INSERT INTO groups (name, created_by)
	VALUES ($1, $2)
	RETURNING id, name, created_by, created_at, updated_at
	`, name, createdBy).Scan(&g.ID, &g.Name, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	// Add members to user_groups table
	for _, member := range members {
		_, err = s.DB.Exec(`
		INSERT INTO user_groups (group_id, user_id)
		VALUES ($1, $2)
		`, g.ID, member)
		if err != nil {
			return nil, err
		}
	}
	return &GroupOutput{
		ID:        g.ID,
		Name:      g.Name,
		CreatedBy: g.CreatedBy,
		Members:   members,
	}, nil
}
