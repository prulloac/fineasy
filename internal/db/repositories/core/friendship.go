package repositories

import (
	"encoding/json"
	"fmt"

	"github.com/prulloac/fineasy/pkg"
)

type Friendship struct {
	pkg.Model
	UserID       uint                    `json:"user_id" validate:"required,min=1"`
	FriendID     uint                    `json:"friend_id" validate:"required,min=1"`
	Status       pkg.SocialRequestStatus `json:"status" validate:"numeric"`
	RelationType pkg.FriendRelationType  `json:"relation_type" validate:"numeric"`
}

func (f *Friendship) String() string {
	out, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf("%+v", f.ID)
	}
	return string(out)
}

func (s *CoreRepository) CreateFriendship(userID, friendID uint) (*Friendship, error) {
	var f Friendship
	err := s.Persistence.QueryRow(`
	INSERT INTO friendships (user_id, friend_id, status, relation_type)
	VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, friend_id, status, relation_type, created_at, updated_at
	`, userID, friendID, pkg.Pending, pkg.Contact).Scan(&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.RelationType, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *CoreRepository) GetFriendshipsByUserID(userID uint) ([]Friendship, error) {
	friends := []Friendship{}
	rows, err := s.Persistence.Query(`
	SELECT id, user_id, friend_id, status, relation_type, created_at, updated_at
	FROM friendships
	WHERE (user_id = $1 OR friend_id = $1) AND status = $2
	`, userID, pkg.Accepted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var f Friendship
		if err := rows.Scan(&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.RelationType, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		friends = append(friends, f)
	}
	return friends, nil
}

func (s *CoreRepository) GetFriendshipByFriendIDAndUserID(fid, uid uint) (*Friendship, error) {
	var f Friendship
	err := s.Persistence.QueryRow(`
	SELECT id, user_id, friend_id, status, relation_type, created_at, updated_at
	FROM friendships
	WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)
	`, uid, fid).Scan(&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.RelationType, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *CoreRepository) GetPendingFriendshipsByUserID(userID uint) ([]Friendship, error) {
	requests := []Friendship{}
	rows, err := s.Persistence.Query(`
	SELECT id, user_id, friend_id, status, relation_type, created_at, updated_at
	FROM friendships
	WHERE (user_id = $1 OR friend_id = $1) AND status = $2
	`, userID, pkg.Pending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Friendship
		if err := rows.Scan(&r.ID, &r.UserID, &r.FriendID, &r.Status, &r.RelationType, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}
	return requests, nil
}

func (s *CoreRepository) AcceptFriendship(userID, friendID uint) (*Friendship, error) {
	f := []Friendship{}
	rows, err := s.Persistence.Query(`
	SELECT id, user_id, friend_id, status, relation_type, created_at, updated_at
	FROM friendships
	WHERE user_id = $1 AND friend_id = $2
	OR user_id = $2 AND friend_id = $1
	`, userID, friendID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Friendship
		if err := rows.Scan(&r.ID, &r.UserID, &r.FriendID, &r.Status, &r.RelationType, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		f = append(f, r)
	}

	for _, r := range f {
		err = s.Persistence.QueryRow(`
		UPDATE friendships
		SET status = $1
		WHERE id = $2
		RETURNING id, user_id, friend_id, status, relation_type, created_at, updated_at
		`, pkg.Accepted, r.ID).Scan(&r.ID, &r.UserID, &r.FriendID, &r.Status, &r.RelationType, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}
	return s.GetFriendshipByFriendIDAndUserID(friendID, userID)
}

func (s *CoreRepository) RejectFriendship(userID, friendID uint) error {
	err := s.Persistence.QueryRow(`
	UPDATE friendships
	SET status = $1
	WHERE (user_id = $2 AND friend_id = $3) OR (user_id = $3 AND friend_id = $2)
	RETURNING id
	`, pkg.Declined, userID, friendID).Scan()
	return err
}
