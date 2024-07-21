package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/prulloac/fineasy/pkg"
)

type UserGroup struct {
	ID       uint                    `json:"id" validate:"required,min=1"`
	UserID   uint                    `json:"user_id" validate:"required,min=1"`
	GroupID  uint                    `json:"group_id" validate:"required,min=1"`
	JoinedAt time.Time               `json:"joined_at" validate:"required,past_time"`
	LeftAt   sql.NullTime            `json:"left_at" validate:"past_time"`
	Status   pkg.SocialRequestStatus `json:"status" validate:"required"`
}

func (ug *UserGroup) String() string {
	out, err := json.Marshal(ug)
	if err != nil {
		return fmt.Sprintf("%+v", ug.ID)
	}
	return string(out)
}

func (s *CoreRepository) InsertUserGroup(userID, groupID uint, status pkg.SocialRequestStatus) (*UserGroup, error) {
	var ug UserGroup
	err := s.Persistence.QueryRow(`
	INSERT INTO user_groups (user_id, group_id, status)
	VALUES ($1, $2, $3)
	RETURNING id, user_id, group_id, joined_at, left_at, status
	`, userID, groupID, status).Scan(&ug.ID, &ug.UserID, &ug.GroupID, &ug.JoinedAt, &ug.LeftAt, &ug.Status)
	if err != nil {
		return nil, err
	}
	return &ug, nil
}

func (s *CoreRepository) GetUserGroup(userGroupID uint) (*UserGroup, error) {
	var ug UserGroup
	err := s.Persistence.QueryRow(`
	SELECT id, user_id, group_id, joined_at, left_at, status
	FROM user_groups
	WHERE id = $1
	`, userGroupID).Scan(&ug.ID, &ug.UserID, &ug.GroupID, &ug.JoinedAt, &ug.LeftAt, &ug.Status)
	if err != nil {
		return nil, err
	}
	return &ug, nil
}

func (s *CoreRepository) GetUserGroupsByUserID(userID uint) ([]UserGroup, error) {
	rows, err := s.Persistence.Query(`
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

func (s *CoreRepository) UpdateUserGroup(userID, groupID uint, status pkg.SocialRequestStatus) (*UserGroup, error) {
	var ug UserGroup
	err := s.Persistence.QueryRow(`
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

func (s *CoreRepository) LeaveGroup(groupID, userID uint) (*UserGroup, error) {
	return s.UpdateUserGroup(userID, groupID, pkg.Left)
}

func (s *CoreRepository) InviteGroupRequest(groupID, userID uint) (*UserGroup, error) {
	return s.InsertUserGroup(userID, groupID, pkg.Invited)
}

func (s *CoreRepository) AcceptGroupRequest(groupID, userID uint) (*UserGroup, error) {
	return s.UpdateUserGroup(userID, groupID, pkg.Accepted)
}

func (s *CoreRepository) RejectGroupRequest(groupID, userID uint) (*UserGroup, error) {
	return s.UpdateUserGroup(userID, groupID, pkg.Declined)
}

func (s *CoreRepository) GetMembershipsByGroupID(groupID uint) ([]UserGroup, error) {
	rows, err := s.Persistence.Query(`
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
