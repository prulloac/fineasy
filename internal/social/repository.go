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
	return s.Persistence.ORM().AutoMigrate(&Group{}, &UserGroup{}, &Friendship{})
}

func (s *SocialRepository) DropTables() error {
	return s.Persistence.ORM().Migrator().DropTable(&Group{}, &UserGroup{}, &Friendship{})
}

func (s *SocialRepository) CreateFriendship(userID, friendID uint) (*Friendship, error) {
	var f Friendship
	err := s.Persistence.ORM().Model(&Friendship{}).Create(&Friendship{
		UserID:       userID,
		FriendID:     friendID,
		Status:       pkg.Pending,
		RelationType: pkg.Contact,
	}).Scan(&f).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *SocialRepository) GetFriendshipsByUserID(userID uint) ([]Friendship, error) {
	friends := []Friendship{}
	rows, err := s.Persistence.ORM().Model(&Friendship{}).
		Where("(user_id = ? OR friend_id = ?) AND status = ?", userID, userID, pkg.Accepted).
		Find(&friends).Rows()
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

func (s *SocialRepository) GetFriendshipByFriendIDAndUserID(fid, uid uint) (*Friendship, error) {
	var f Friendship
	err := s.Persistence.ORM().Model(&Friendship{}).Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", uid, fid, fid, uid).First(&f).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *SocialRepository) GetPendingFriendshipsByUserID(userID uint) ([]Friendship, error) {
	requests := []Friendship{}
	rows, err := s.Persistence.ORM().Debug().Model(&Friendship{}).Find(&requests, "(user_id = ? OR friend_id = ?) AND status = ?", userID, userID, pkg.Pending).Rows()
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

func (s *SocialRepository) AcceptFriendship(userID, friendID uint) (*Friendship, error) {
	f := []Friendship{}
	err := s.Persistence.ORM().Model(&Friendship{}).Debug().
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Or("user_id = ? AND friend_id = ?", friendID, userID).
		Find(&f).Error
	if err != nil {
		return nil, err
	}
	for _, r := range f {
		err = s.Persistence.ORM().Model(&Friendship{}).Where("id = ?", r.ID).Update("status", pkg.Accepted).Error
		if err != nil {
			return nil, err
		}
	}
	return s.GetFriendshipByFriendIDAndUserID(friendID, userID)
}

func (s *SocialRepository) RejectFriendship(userID, friendID uint) error {
	err := s.Persistence.ORM().Model(&Friendship{}).Delete(&Friendship{}, "(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", userID, friendID, friendID, userID).Error
	return err
}

func (s *SocialRepository) CreateGroup(name string, createdBy uint) (*Group, error) {
	var g Group
	err := s.Persistence.ORM().Model(&Group{}).Create(&Group{
		Name:      name,
		CreatedBy: createdBy,
	}).Scan(&g).Error
	if err != nil {
		return nil, err
	}
	s.Persistence.ORM().Model(&Group{}).Save(&g)
	// Add members to user_groups table
	_, err = s.InsertUserGroup(createdBy, g.ID, pkg.Accepted)
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

func (s *SocialRepository) InsertUserGroup(userID, groupID uint, status pkg.SocialRequestStatus) (*UserGroup, error) {
	var ug UserGroup
	err := s.Persistence.ORM().Model(&UserGroup{}).Create(&UserGroup{
		UserID:  userID,
		GroupID: groupID,
		Status:  status,
	}).Scan(&ug).Error
	s.Persistence.ORM().Model(&UserGroup{}).Save(&ug)
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

func (s *SocialRepository) UpdateUserGroup(userID, groupID uint, status pkg.SocialRequestStatus) (*UserGroup, error) {
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
		if err = pkg.ValidateStruct(ug); err != nil {
			return nil, err
		}
		memberships = append(memberships, ug)
	}
	return memberships, nil
}
