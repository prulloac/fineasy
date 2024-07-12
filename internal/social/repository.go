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
		member_count INT NOT NULL DEFAULT 0,
		created_by INT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS user_groups (
		id SERIAL PRIMARY KEY,
		group_id INT NOT NULL,
		user_id INT NOT NULL,
		joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		left_at TIMESTAMP,
		status INTEGER NOT NULL
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

func (s *SocialRepository) AddFriend(userID, friendID int) (*FriendRequest, error) {
	var f FriendRequest
	err := s.DB.QueryRow(`
	INSERT INTO friend_requests (user_id, friend_id, status)
	VALUES ($1, $2, $3)
	RETURNING id, user_id, friend_id, status, created_at, updated_at
	`, userID, friendID, Pending).Scan(&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *SocialRepository) GetFriends(userID int) ([]Friend, error) {
	rows, err := s.DB.Query(`
	SELECT user_id, friend_id, relation_type
	FROM friends
	WHERE user_id = $1 OR friend_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	friends := []Friend{}
	for rows.Next() {
		var f Friend
		if err := rows.Scan(&f.UserID, &f.FriendID, &f.RelationType); err != nil {
			return nil, err
		}
		if err = pkg.ValidateStruct(f); err != nil {
			return nil, err
		}
		friends = append(friends, f)
	}
	return friends, nil
}

func (s *SocialRepository) GetFriendRequests(userID int) ([]FriendRequest, error) {
	rows, err := s.DB.Query(`
	SELECT user_id, friend_id, status
	FROM friend_requests
	WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := []FriendRequest{}
	for rows.Next() {
		var f FriendRequest
		if err := rows.Scan(&f.UserID, &f.FriendID, &f.Status); err != nil {
			return nil, err
		}
		if err = pkg.ValidateStruct(f); err != nil {
			return nil, err
		}
		requests = append(requests, f)
	}
	return requests, nil
}

func (s *SocialRepository) AcceptFriendRequest(userID, friendID int) (*FriendRequest, error) {
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
	return &f, nil
}

func (s *SocialRepository) RejectFriendRequest(userID, friendID int) (*FriendRequest, error) {
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
	return &f, nil
}

func (s *SocialRepository) DeleteFriend(userID, friendID int) error {
	_, err := s.DB.Exec(`
	DELETE FROM friends
	WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)
	`, userID, friendID)
	return err
}

func (s *SocialRepository) CreateGroup(name string, createdBy int) (*Group, error) {
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
	_, err = s.InsertUserGroup(createdBy, g.ID, Accepted)
	return &g, nil
}

func (s *SocialRepository) GetGroupByID(groupID int) (*Group, error) {
	var g Group
	err := s.DB.QueryRow(`
	SELECT id, name, created_by, created_at, updated_at
	FROM groups
	WHERE id = $1
	`, groupID).Scan(&g.ID, &g.Name, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *SocialRepository) GetGroupsByUserID(userID int) ([]Group, error) {
	rows, err := s.DB.Query(`
	SELECT g.id, g.name, g.created_by, g.created_at, g.updated_at, g.member_count
	FROM groups g
	JOIN user_groups ug ON g.id = ug.group_id
	WHERE ug.user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []Group{}
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Name, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt, &g.MemberCount); err != nil {
			return nil, err
		}
		if err = pkg.ValidateStruct(g); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func (s *SocialRepository) UpdateGroup(groupID int, name string) (*Group, error) {
	var g Group
	err := s.DB.QueryRow(`
	UPDATE groups
	SET name = $2
	WHERE id = $1
	RETURNING id, name, created_by, created_at, updated_at
	`, groupID, name).Scan(&g.ID, &g.Name, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *SocialRepository) InsertUserGroup(userID, groupID int, status SocialRequestStatus) (*UserGroup, error) {
	var ug UserGroup
	err := s.DB.QueryRow(`
	INSERT INTO user_groups (group_id, user_id, status)
	VALUES ($1, $2, $3)
	RETURNING id, user_id, group_id, joined_at, left_at, status
	`, groupID, userID, status).Scan(&ug.ID, &ug.UserID, &ug.GroupID, &ug.JoinedAt, &ug.LeftAt, &ug.Status)
	if err != nil {
		return nil, err
	}
	return &ug, nil
}

func (s *SocialRepository) GetUserGroup(userGroupID int) (*UserGroup, error) {
	var ug UserGroup
	err := s.DB.QueryRow(`
	SELECT id, user_id, group_id, joined_at, left_at, status
	FROM user_groups
	WHERE id = $1
	`, userGroupID).Scan(&ug.ID, &ug.UserID, &ug.GroupID, &ug.JoinedAt, &ug.LeftAt, &ug.Status)
	if err != nil {
		return nil, err
	}
	return &ug, nil
}

func (s *SocialRepository) GetUserGroupsByUserID(userID int) ([]UserGroup, error) {
	rows, err := s.DB.Query(`
	SELECT id, user_id, group_id, joined_at, left_at, status
	FROM user_groups
	WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userGroups := []UserGroup{}
	for rows.Next() {
		var ug UserGroup
		if err := rows.Scan(&ug.ID, &ug.UserID, &ug.GroupID, &ug.JoinedAt, &ug.LeftAt, &ug.Status); err != nil {
			return nil, err
		}
		if err = pkg.ValidateStruct(ug); err != nil {
			return nil, err
		}
		userGroups = append(userGroups, ug)
	}
	return userGroups, nil
}

func (s *SocialRepository) UpdateUserGroup(userID, groupID int, status SocialRequestStatus) (*UserGroup, error) {
	var ug UserGroup
	err := s.DB.QueryRow(`
	UPDATE user_groups
	SET status = $3
	WHERE user_id = $1 AND group_id = $2
	RETURNING id, user_id, group_id, joined_at, left_at, status
	`, userID, groupID, status).Scan(&ug.ID, &ug.UserID, &ug.GroupID, &ug.JoinedAt, &ug.LeftAt, &ug.Status)
	if err != nil {
		return nil, err
	}
	return &ug, nil
}

func (s *SocialRepository) LeaveGroup(groupID, userID int) (*UserGroup, error) {
	return s.UpdateUserGroup(userID, groupID, Left)
}

func (s *SocialRepository) InviteGroupRequest(groupID, userID int) (*UserGroup, error) {
	return s.InsertUserGroup(userID, groupID, Invited)
}

func (s *SocialRepository) AcceptGroupRequest(groupID, userID int) (*UserGroup, error) {
	return s.UpdateUserGroup(userID, groupID, Accepted)
}

func (s *SocialRepository) RejectGroupRequest(groupID, userID int) (*UserGroup, error) {
	return s.UpdateUserGroup(userID, groupID, Declined)
}
