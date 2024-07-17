package social

import (
	p "github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/pkg"
)

type SocialRepository struct {
	Persistence *p.Persistence
}

func NewSocialRepository(persistence *p.Persistence) *SocialRepository {
	return &SocialRepository{persistence}
}

func (s *SocialRepository) Close() {
	s.Persistence.Close()
}

func (s *SocialRepository) CreateTables() error {
	_, err := s.Persistence.SQL().Exec(`
	CREATE TABLE IF NOT EXISTS friendships (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		friend_id INT NOT NULL,
		status INTEGER NOT NULL,
		relation_type INTEGER NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP
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

func (s *SocialRepository) DropTables() error {
	_, err := s.Persistence.SQL().Exec(`
	DROP TABLE IF EXISTS user_groups;
	DROP TABLE IF EXISTS groups;
	DROP TABLE IF EXISTS friendships;
	`)
	return err
}

func (s *SocialRepository) CreateFriendship(userID, friendID uint) (*Friendship, error) {
	var f Friendship
	err := s.Persistence.SQL().QueryRow(`
	INSERT INTO friendships (user_id, friend_id, status, relation_type)
	VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, friend_id, status, relation_type, created_at, updated_at
	`, userID, friendID, pkg.Pending, pkg.Contact).Scan(&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.RelationType, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *SocialRepository) GetFriendshipsByUserID(userID uint) ([]Friendship, error) {
	friends := []Friendship{}
	rows, err := s.Persistence.SQL().Query(`
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

func (s *SocialRepository) GetFriendshipByFriendIDAndUserID(fid, uid uint) (*Friendship, error) {
	var f Friendship
	err := s.Persistence.SQL().QueryRow(`
	SELECT id, user_id, friend_id, status, relation_type, created_at, updated_at
	FROM friendships
	WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)
	`, uid, fid).Scan(&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.RelationType, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *SocialRepository) GetPendingFriendshipsByUserID(userID uint) ([]Friendship, error) {
	requests := []Friendship{}
	rows, err := s.Persistence.SQL().Query(`
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

func (s *SocialRepository) AcceptFriendship(userID, friendID uint) (*Friendship, error) {
	f := []Friendship{}
	rows, err := s.Persistence.SQL().Query(`
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
		err = s.Persistence.SQL().QueryRow(`
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

func (s *SocialRepository) RejectFriendship(userID, friendID uint) error {
	err := s.Persistence.SQL().QueryRow(`
	UPDATE friendships
	SET status = $1
	WHERE (user_id = $2 AND friend_id = $3) OR (user_id = $3 AND friend_id = $2)
	RETURNING id
	`, pkg.Declined, userID, friendID).Scan()
	return err
}

func (s *SocialRepository) CreateGroup(name string, createdBy uint) (*Group, error) {
	var g Group
	err := s.Persistence.SQL().QueryRow(`
	INSERT INTO groups (name, created_by)
	VALUES ($1, $2)
	RETURNING id, name, created_by, created_at, updated_at
	`, name, createdBy).Scan(&g.ID, &g.Name, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_, err = s.InsertUserGroup(createdBy, g.ID, pkg.Accepted)
	if err != nil {
		return nil, err
	}
	g.MemberCount = 1
	return &g, nil
}

func (s *SocialRepository) GetGroupByID(groupID uint) (*Group, error) {
	var g Group
	err := s.Persistence.SQL().QueryRow(`
	SELECT id, name, created_by, created_at, updated_at
	FROM groups
	WHERE id = $1
	`, groupID).Scan(&g.ID, &g.Name, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *SocialRepository) GetGroupByUserID(gid, uid uint) (*Group, error) {
	var g Group
	err := s.Persistence.SQL().QueryRow(`
	SELECT g.id, g.name, g.created_by, g.created_at, g.updated_at
	FROM groups g
	JOIN user_groups ug ON g.id = ug.group_id
	WHERE g.id = $1 AND ug.user_id = $2
	`, gid, uid).Scan(&g.ID, &g.Name, &g.CreatedBy, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *SocialRepository) GetGroupsByUserID(userID uint) ([]Group, error) {
	rows, err := s.Persistence.SQL().Query(`
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
		groups = append(groups, g)
	}
	return groups, nil
}

func (s *SocialRepository) UpdateGroup(groupID uint, name string) (*Group, error) {
	var g Group
	err := s.Persistence.SQL().QueryRow(`
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

func (s *SocialRepository) InsertUserGroup(userID, groupID uint, status pkg.SocialRequestStatus) (*UserGroup, error) {
	var ug UserGroup
	err := s.Persistence.SQL().QueryRow(`
	INSERT INTO user_groups (user_id, group_id, status)
	VALUES ($1, $2, $3)
	RETURNING id, user_id, group_id, joined_at, left_at, status
	`, userID, groupID, status).Scan(&ug.ID, &ug.UserID, &ug.GroupID, &ug.JoinedAt, &ug.LeftAt, &ug.Status)
	if err != nil {
		return nil, err
	}
	return &ug, nil
}

func (s *SocialRepository) GetUserGroup(userGroupID uint) (*UserGroup, error) {
	var ug UserGroup
	err := s.Persistence.SQL().QueryRow(`
	SELECT id, user_id, group_id, joined_at, left_at, status
	FROM user_groups
	WHERE id = $1
	`, userGroupID).Scan(&ug.ID, &ug.UserID, &ug.GroupID, &ug.JoinedAt, &ug.LeftAt, &ug.Status)
	if err != nil {
		return nil, err
	}
	return &ug, nil
}

func (s *SocialRepository) GetUserGroupsByUserID(userID uint) ([]UserGroup, error) {
	rows, err := s.Persistence.SQL().Query(`
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
		userGroups = append(userGroups, ug)
	}
	return userGroups, nil
}

func (s *SocialRepository) UpdateUserGroup(userID, groupID uint, status pkg.SocialRequestStatus) (*UserGroup, error) {
	var ug UserGroup
	err := s.Persistence.SQL().QueryRow(`
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

func (s *SocialRepository) LeaveGroup(groupID, userID uint) (*UserGroup, error) {
	return s.UpdateUserGroup(userID, groupID, pkg.Left)
}

func (s *SocialRepository) InviteGroupRequest(groupID, userID uint) (*UserGroup, error) {
	return s.InsertUserGroup(userID, groupID, pkg.Invited)
}

func (s *SocialRepository) AcceptGroupRequest(groupID, userID uint) (*UserGroup, error) {
	return s.UpdateUserGroup(userID, groupID, pkg.Accepted)
}

func (s *SocialRepository) RejectGroupRequest(groupID, userID uint) (*UserGroup, error) {
	return s.UpdateUserGroup(userID, groupID, pkg.Declined)
}

func (s *SocialRepository) GetMembershipsByGroupID(groupID uint) ([]UserGroup, error) {
	rows, err := s.Persistence.SQL().Query(`
	SELECT id, user_id, group_id, joined_at, left_at, status
	FROM user_groups
	WHERE group_id = $1
	`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memberships := []UserGroup{}
	for rows.Next() {
		var ug UserGroup
		if err := rows.Scan(&ug.ID, &ug.UserID, &ug.GroupID, &ug.JoinedAt, &ug.LeftAt, &ug.Status); err != nil {
			return nil, err
		}
		memberships = append(memberships, ug)
	}
	return memberships, nil
}
