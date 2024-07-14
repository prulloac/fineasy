package social

import (
	"log"

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

func (s *SocialRepository) DropTables() error {
	_, err := s.Persistence.SQL().Exec(`
	DROP TABLE IF EXISTS user_groups;
	DROP TABLE IF EXISTS groups;
	DROP TABLE IF EXISTS friend_requests;
	DROP TABLE IF EXISTS friends;
	`)
	return err
}

func (s *SocialRepository) AddFriend(userID, friendID uint) (*FriendRequest, error) {
	var f FriendRequest
	err := s.Persistence.SQL().QueryRow(`
	INSERT INTO friend_requests (user_id, friend_id, status)
	VALUES ($1, $2, $3)
	RETURNING id, user_id, friend_id, status, created_at, updated_at
	`, userID, friendID, Pending).Scan(&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *SocialRepository) GetFriends(userID uint) ([]Friend, error) {
	friends := []Friend{}

	rows, err := s.Persistence.ORM().Debug().
		Table("friends").Find(&friends).Where("user_id = ?", userID).Or("friend_id = ?", userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for _, f := range friends {
		if err = pkg.ValidateStruct(f); err != nil {
			return nil, err
		}
	}
	return friends, nil
}

func (s *SocialRepository) GetFriend(fid, uid uint) (*Friend, error) {
	var f Friend
	err := s.Persistence.SQL().QueryRow(`
	SELECT user_id, friend_id, relation_type
	FROM friends
	WHERE user_id = $1 AND friend_id = $2
	`, uid, fid).Scan(&f.UserID, &f.FriendID, &f.RelationType)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *SocialRepository) GetFriendRequests(userID uint) ([]FriendRequest, error) {
	requests := []FriendRequest{}
	rows, err := s.Persistence.ORM().Debug().Table("friend_requests").Find(&requests, "user_id = ?", userID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for _, r := range requests {
		if err = pkg.ValidateStruct(r); err != nil {
			return nil, err
		}
	}
	return requests, nil
}

func (s *SocialRepository) AcceptFriendRequest(userID, friendID uint) (*FriendRequest, error) {
	var f FriendRequest
	err := s.Persistence.SQL().QueryRow(`
	UPDATE friend_requests
	SET status = $3
	WHERE user_id = $1 AND friend_id = $2
	RETURNING id, user_id, friend_id, status, created_at, updated_at
	`, userID, friendID, Accepted).Scan(&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	// both users are now friends
	_, err = s.Persistence.SQL().Exec(`
	INSERT INTO friends (user_id, friend_id, relation_type)
	VALUES ($1, $2, $3)
	`, userID, friendID, Contact)
	if err != nil {
		return nil, err
	}
	_, err = s.Persistence.SQL().Exec(`
	INSERT INTO friends (user_id, friend_id, relation_type)
	VALUES ($1, $2, $3)
	`, friendID, userID, Contact)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *SocialRepository) RejectFriendRequest(userID, friendID uint) (*FriendRequest, error) {
	var f FriendRequest
	err := s.Persistence.SQL().QueryRow(`
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

func (s *SocialRepository) DeleteFriend(userID, friendID uint) error {
	// Delete friendship from friends table
	_, err := s.Persistence.SQL().Exec(`
	DELETE FROM friends
	WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)
	`, userID, friendID)
	if err != nil {
		return err
	}
	_, err = s.Persistence.SQL().Exec(`
	DELETE FROM friends
	WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)
	`, friendID, userID)
	if err != nil {
		return err
	}
	// Delete friend request from friend_requests table
	_, err = s.Persistence.SQL().Exec(`
	DELETE FROM friend_requests
	WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)
	`, userID, friendID)
	if err != nil {
		return err
	}
	return nil
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
	// Add members to user_groups table
	_, err = s.InsertUserGroup(createdBy, g.ID, Accepted)
	if err != nil {
		return nil, err
	}
	g.MemberCount = 1
	return &g, nil
}

func (s *SocialRepository) GetGroupByID(groupID uint) (*Group, error) {
	var g Group
	err := s.Persistence.ORM().Table("groups").Where("id = ?", groupID).First(&g).Error
	if err != nil {
		return nil, err
	}
	log.Printf("🔍 Found group: %v", g)
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
		if err = pkg.ValidateStruct(g); err != nil {
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

func (s *SocialRepository) InsertUserGroup(userID, groupID uint, status SocialRequestStatus) (*UserGroup, error) {
	var ug UserGroup
	err := s.Persistence.SQL().QueryRow(`
	INSERT INTO user_groups (group_id, user_id, status)
	VALUES ($1, $2, $3)
	RETURNING id, user_id, group_id, joined_at, left_at, status
	`, groupID, userID, status).Scan(&ug.ID, &ug.UserID, &ug.GroupID, &ug.JoinedAt, &ug.LeftAt, &ug.Status)
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
		if err = pkg.ValidateStruct(ug); err != nil {
			return nil, err
		}
		userGroups = append(userGroups, ug)
	}
	return userGroups, nil
}

func (s *SocialRepository) UpdateUserGroup(userID, groupID uint, status SocialRequestStatus) (*UserGroup, error) {
	var ug UserGroup
	err := s.Persistence.ORM().Table("user_groups").Where("user_id = ? AND group_id = ?", userID, groupID).First(&ug).Error
	if err != nil {
		return nil, err
	}
	ug.Status = status
	err = s.Persistence.ORM().Table("user_groups").Save(&ug).Error
	if err != nil {
		return nil, err
	}
	return &ug, nil
}

func (s *SocialRepository) LeaveGroup(groupID, userID uint) (*UserGroup, error) {
	return s.UpdateUserGroup(userID, groupID, Left)
}

func (s *SocialRepository) InviteGroupRequest(groupID, userID uint) (*UserGroup, error) {
	return s.InsertUserGroup(userID, groupID, Invited)
}

func (s *SocialRepository) AcceptGroupRequest(groupID, userID uint) (*UserGroup, error) {
	return s.UpdateUserGroup(userID, groupID, Accepted)
}

func (s *SocialRepository) RejectGroupRequest(groupID, userID uint) (*UserGroup, error) {
	return s.UpdateUserGroup(userID, groupID, Declined)
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
		if err = pkg.ValidateStruct(ug); err != nil {
			return nil, err
		}
		memberships = append(memberships, ug)
	}
	return memberships, nil
}
