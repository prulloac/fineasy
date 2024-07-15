package social

import (
	"log"
	"slices"

	e "github.com/prulloac/fineasy/internal/errors"
	p "github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/pkg"
)

type Service struct {
	repo *SocialRepository
}

func NewService(per *p.Persistence) *Service {
	instance := &Service{}
	instance.repo = NewSocialRepository(per)
	return instance
}

func (s *Service) Close() {
	s.repo.Close()
}

func (s *Service) AddFriendship(fid, uid uint) (*FriendRequestOutput, error) {
	fr, err := s.repo.CreateFriendship(uid, fid)
	if err != nil {
		return nil, err
	}
	out := &FriendRequestOutput{
		UserID:   fr.UserID,
		FriendID: fr.FriendID,
		Status:   fr.Status.String(),
	}
	if err = pkg.ValidateStruct(out); err != nil {
		log.Printf("⚠️ Error adding friend: %s", err)
		return nil, err
	}
	return out, nil
}

func (s *Service) GetFriendships(uid uint) ([]FriendShipOutput, error) {
	fs, err := s.repo.GetFriendshipsByUserID(uid)
	if err != nil {
		return nil, err
	}
	out := []FriendShipOutput{}
	for _, f := range fs {
		e := FriendShipOutput{
			UserID:       f.UserID,
			FriendID:     f.FriendID,
			RelationType: f.RelationType.String(),
		}
		if err = pkg.ValidateStruct(e); err != nil {
			log.Printf("⚠️ Error getting friends: %s", err)
			return nil, err
		}
		out = append(out, e)
	}
	return out, nil
}

func (s *Service) GetFriendship(fid, uid uint) (*FriendShipOutput, error) {
	f, err := s.repo.GetFriendshipByFriendIDAndUserID(fid, uid)
	if err != nil {
		return nil, err
	}
	out := &FriendShipOutput{
		UserID:       f.UserID,
		FriendID:     f.FriendID,
		RelationType: f.RelationType.String(),
	}
	if err = pkg.ValidateStruct(out); err != nil {
		log.Printf("⚠️ Error getting friend: %s", err)
		return nil, err
	}
	return out, nil
}

func (s *Service) GetPendingFriendships(uid uint) ([]FriendRequestOutput, error) {
	frs, err := s.repo.GetPendingFriendshipsByUserID(uid)
	if err != nil {
		return nil, err
	}
	out := []FriendRequestOutput{}
	for _, fr := range frs {
		e := FriendRequestOutput{
			UserID:   fr.UserID,
			FriendID: fr.FriendID,
			Status:   fr.Status.String(),
		}
		if err = pkg.ValidateStruct(e); err != nil {
			log.Printf("⚠️ Error getting friend requests: %s", err)
			return nil, err
		}
		out = append(out, e)
	}
	return out, nil
}

func (s *Service) AcceptFriendship(status string, fid, uid uint) (*FriendRequestOutput, error) {
	var fr *Friendship
	var err error
	if status == "Accepted" {
		fr, err = s.repo.AcceptFriendship(uid, fid)
	}
	if err != nil {
		return nil, err
	}
	out := &FriendRequestOutput{
		UserID:   fr.UserID,
		FriendID: fr.FriendID,
		Status:   fr.Status.String(),
	}
	if err = pkg.ValidateStruct(out); err != nil {
		log.Printf("⚠️ Error updating friend request: %s", err)
		return nil, err
	}
	return out, nil
}

func (s *Service) RejectFriendship(fid, uid uint) ([]FriendShipOutput, error) {
	err := s.repo.RejectFriendship(uid, fid)
	if err != nil {
		return nil, err
	}
	return s.GetFriendships(uid)
}

func (s *Service) CreateGroup(name string, uid uint) (*GroupBriefOutput, error) {
	g, err := s.repo.CreateGroup(name, uid)
	if err != nil {
		return nil, err
	}
	out := &GroupBriefOutput{
		ID:          g.ID,
		Name:        g.Name,
		MemberCount: g.MemberCount,
		CreatedBy:   g.CreatedBy,
	}
	if err = pkg.ValidateStruct(out); err != nil {
		log.Printf("⚠️ Error creating group: %s", err)
		return nil, err
	}
	return out, nil
}

func (s *Service) GetGroupByID(gid, uid uint) (*GroupFullOutput, error) {
	g, err := s.repo.GetGroupByUserID(gid, uid)
	if err != nil {
		return nil, err
	}
	ms, err := s.repo.GetMembershipsByGroupID(gid)
	if err != nil {
		return nil, err
	}
	out := &GroupFullOutput{
		GroupID:     g.ID,
		Name:        g.Name,
		MemberCount: g.MemberCount,
		CreatedBy:   g.CreatedBy,
		Memberships: []MembershipOutput{},
	}
	for _, m := range ms {
		e := MembershipOutput{
			UserID:   m.UserID,
			Status:   m.Status.String(),
			JoinedAt: m.JoinedAt.String(),
		}
		if err = pkg.ValidateStruct(e); err != nil {
			log.Printf("⚠️ Error getting group by id: %s", err)
			return nil, err
		}
		out.Memberships = append(out.Memberships, e)
	}
	if err = pkg.ValidateStruct(out); err != nil {
		log.Printf("⚠️ Error getting group by id: %s", err)
		return nil, err
	}
	return out, nil
}

func (s *Service) GetGroup(id uint) (*GroupBriefOutput, error) {
	g, err := s.repo.GetGroupByID(id)
	if err != nil {
		return nil, err
	}
	out := &GroupBriefOutput{
		ID:          g.ID,
		Name:        g.Name,
		MemberCount: g.MemberCount,
		CreatedBy:   g.CreatedBy,
	}
	if err = pkg.ValidateStruct(out); err != nil {
		log.Printf("⚠️ Error getting group: %s", err)
		return nil, err
	}
	return out, nil
}

func (s *Service) GetUserGroups(uid uint) ([]UserGroupOutput, error) {
	ugs, err := s.repo.GetUserGroupsByUserID(uid)
	if err != nil {
		return nil, err
	}
	out := []UserGroupOutput{}
	for _, ug := range ugs {
		g, err := s.repo.GetGroupByID(ug.GroupID)
		if err != nil {
			return nil, err
		}
		leftAt := ""
		if ug.LeftAt.Valid {
			leftAt = ug.LeftAt.Time.String()
		}
		e := UserGroupOutput{
			UserID:      ug.UserID,
			GroupID:     ug.GroupID,
			MemberCount: g.MemberCount,
			Status:      ug.Status.String(),
			Group:       g.Name,
			CreatedBy:   g.CreatedBy,
			JoinedAt:    ug.JoinedAt.String(),
			LeftAt:      leftAt,
		}
		if err = pkg.ValidateStruct(e); err != nil {
			log.Printf("⚠️ Error getting user groups: %s", err)
			return nil, err
		}
		out = append(out, e)
	}
	return out, nil
}

func (s *Service) UpdateGroup(name string, gid, uid uint) (*GroupBriefOutput, error) {
	groups, err := s.repo.GetGroupsByUserID(uid)
	if err != nil {
		return nil, err
	}
	userIsInGroup := slices.ContainsFunc(groups, func(g Group) bool {
		return g.ID == gid
	})

	if !userIsInGroup {
		err := &e.ErrForbidden{}
		log.Printf("⚠️ Error updating group: %s", err)
		return nil, err
	}

	g, err := s.repo.UpdateGroup(gid, name)
	if err != nil {
		return nil, err
	}
	out := &GroupBriefOutput{
		ID:          g.ID,
		Name:        g.Name,
		MemberCount: g.MemberCount,
		CreatedBy:   g.CreatedBy,
	}
	if err = pkg.ValidateStruct(out); err != nil {
		log.Printf("⚠️ Error updating group: %s", err)
		return nil, err
	}
	return out, nil
}

func (s *Service) UpdateUserGroup(status string, gid, uid uint) (*UserGroupOutput, error) {
	var ug *UserGroup
	var err error
	switch status {
	case "Accepted":
		ug, err = s.repo.AcceptGroupRequest(gid, uid)
	case "Rejected":
		ug, err = s.repo.RejectGroupRequest(gid, uid)
	case "Left":
		ug, err = s.repo.LeaveGroup(gid, uid)
	default:
		err := &e.ErrBadRequest{}
		log.Printf("⚠️ Error joining group: %s", err)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	g, err := s.GetGroup(ug.GroupID)
	if err != nil {
		return nil, err
	}
	leftAt := ""
	if ug.LeftAt.Valid {
		leftAt = ug.LeftAt.Time.String()
	}
	out := &UserGroupOutput{
		UserID:      ug.UserID,
		GroupID:     ug.GroupID,
		MemberCount: g.MemberCount,
		Status:      ug.Status.String(),
		Group:       g.Name,
		JoinedAt:    ug.JoinedAt.String(),
		LeftAt:      leftAt,
		CreatedBy:   g.CreatedBy,
	}
	if err = pkg.ValidateStruct(out); err != nil {
		log.Printf("⚠️ Error joining group: %s", err)
		return nil, err
	}
	return out, nil
}

func (s *Service) InviteUserGroup(gid, uid uint) (*GroupFullOutput, error) {
	ug, err := s.repo.InviteGroupRequest(gid, uid)
	if err != nil {
		return nil, err
	}
	_, err = s.GetGroup(ug.GroupID)
	if err != nil {
		return nil, err
	}
	return s.GetGroupFullData(ug.GroupID)
}

func (s *Service) GetGroupFullData(gid uint) (*GroupFullOutput, error) {
	g, err := s.repo.GetGroupByID(gid)
	if err != nil {
		return nil, err
	}
	ms, err := s.repo.GetMembershipsByGroupID(gid)
	if err != nil {
		return nil, err
	}
	out := &GroupFullOutput{
		GroupID:     g.ID,
		Name:        g.Name,
		MemberCount: g.MemberCount,
		CreatedBy:   g.CreatedBy,
		Memberships: []MembershipOutput{},
	}
	for _, m := range ms {
		e := MembershipOutput{
			UserID:   m.UserID,
			Status:   m.Status.String(),
			JoinedAt: m.JoinedAt.String(),
		}
		if err = pkg.ValidateStruct(e); err != nil {
			log.Printf("⚠️ Error getting group full data: %s", err)
			return nil, err
		}
		out.Memberships = append(out.Memberships, e)
	}
	if err = pkg.ValidateStruct(out); err != nil {
		log.Printf("⚠️ Error getting group full data: %s", err)
		return nil, err
	}
	return out, nil
}
