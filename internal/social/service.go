package social

import (
	"log"
	"slices"

	e "github.com/prulloac/fineasy/internal/errors"
	p "github.com/prulloac/fineasy/internal/persistence"
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
	return &FriendRequestOutput{
		UserID:   fr.UserID,
		FriendID: fr.FriendID,
		Status:   fr.Status.String(),
	}, nil
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
		out = append(out, e)
	}
	return out, nil
}

func (s *Service) GetFriendship(fid, uid uint) (*FriendShipOutput, error) {
	f, err := s.repo.GetFriendshipByFriendIDAndUserID(fid, uid)
	if err != nil {
		return nil, err
	}
	return &FriendShipOutput{
		UserID:       f.UserID,
		FriendID:     f.FriendID,
		RelationType: f.RelationType.String(),
	}, nil
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
	return &FriendRequestOutput{
		UserID:   fr.UserID,
		FriendID: fr.FriendID,
		Status:   fr.Status.String(),
	}, nil
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
	return &GroupBriefOutput{
		ID:          g.ID,
		Name:        g.Name,
		MemberCount: g.MemberCount,
		CreatedBy:   g.CreatedBy,
	}, nil
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
		out.Memberships = append(out.Memberships, e)
	}
	return out, nil
}

func (s *Service) GetGroup(id uint) (*GroupBriefOutput, error) {
	g, err := s.repo.GetGroupByID(id)
	if err != nil {
		return nil, err
	}
	return &GroupBriefOutput{
		ID:          g.ID,
		Name:        g.Name,
		MemberCount: g.MemberCount,
		CreatedBy:   g.CreatedBy,
	}, nil
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
	return &GroupBriefOutput{
		ID:          g.ID,
		Name:        g.Name,
		MemberCount: g.MemberCount,
		CreatedBy:   g.CreatedBy,
	}, nil
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
	return &UserGroupOutput{
		UserID:      ug.UserID,
		GroupID:     ug.GroupID,
		MemberCount: g.MemberCount,
		Status:      ug.Status.String(),
		Group:       g.Name,
		JoinedAt:    ug.JoinedAt.String(),
		LeftAt:      leftAt,
		CreatedBy:   g.CreatedBy,
	}, nil
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
		out.Memberships = append(out.Memberships, e)
	}
	return out, nil
}
