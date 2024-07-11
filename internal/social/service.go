package social

import (
	"log"
	"slices"

	"github.com/golang-jwt/jwt/v5"
	e "github.com/prulloac/fineasy/internal/errors"
	"github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/pkg"
)

type Service struct {
	repo        *SocialRepository
	persistence *persistence.Persistence
}

func NewService() *Service {
	instance := &Service{}
	instance.persistence = persistence.NewConnection()
	instance.repo = NewSocialRepository(instance.persistence.Session())
	return instance
}

func (s *Service) Close() {
	s.persistence.Close()
}

func (s *Service) AddFriend(i AddFriendInput, t *jwt.Token) (*FriendRequestOutput, error) {
	uid := t.Claims.(jwt.MapClaims)["uid"].(float64)
	if int(uid) != i.UserID {
		err := &e.ErrForbidden{}
		log.Printf("⚠️ Error adding friend: %s", err)
		return nil, err
	}
	if err := pkg.ValidateStruct(i); err != nil {
		log.Printf("⚠️ Error adding friend: %s", err)
		return nil, err
	}
	fr, err := s.repo.AddFriend(i.UserID, i.FriendID)
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

func (s *Service) GetFriends(t *jwt.Token) ([]FriendShipOutput, error) {
	uid := t.Claims.(jwt.MapClaims)["uid"].(float64)
	fs, err := s.repo.GetFriends(int(uid))
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

func (s *Service) GetFriendRequests(t *jwt.Token) ([]FriendRequestOutput, error) {
	uid := t.Claims.(jwt.MapClaims)["uid"].(float64)
	frs, err := s.repo.GetFriendRequests(int(uid))
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

func (s *Service) UpdateFriendRequest(i UpdateFriendRequestInput, t *jwt.Token) (*FriendRequestOutput, error) {
	uid := t.Claims.(jwt.MapClaims)["uid"].(float64)
	if int(uid) != i.UserID {
		err := &e.ErrForbidden{}
		log.Printf("⚠️ Error updating friend request: %s", err)
		return nil, err
	}
	if i.Status != "Accepted" && i.Status != "Rejected" {
		err := &e.ErrBadRequest{}
		log.Printf("⚠️ Error updating friend request: %s", err)
		return nil, err
	}
	if err := pkg.ValidateStruct(i); err != nil {
		log.Printf("⚠️ Error updating friend request: %s", err)
		return nil, err
	}
	var fr *FriendRequest
	var err error
	if i.Status == "Accepted" {
		fr, err = s.repo.AcceptFriendRequest(i.UserID, i.FriendID)
	} else {
		fr, err = s.repo.RejectFriendRequest(i.UserID, i.FriendID)
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

func (s *Service) CreateGroup(i CreateGroupInput, t *jwt.Token) (*GroupBriefOutput, error) {
	uid := t.Claims.(jwt.MapClaims)["uid"].(float64)
	if err := pkg.ValidateStruct(i); err != nil {
		log.Printf("⚠️ Error creating group: %s", err)
		return nil, err
	}
	g, err := s.repo.CreateGroup(i.Name, int(uid))
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

func (s *Service) GetGroup(id int) (*GroupBriefOutput, error) {
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

func (s *Service) GetUserGroups(t *jwt.Token) ([]UserGroupOutput, error) {
	uid := t.Claims.(jwt.MapClaims)["uid"].(float64)
	ugs, err := s.repo.GetUserGroupsByUserID(int(uid))
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

func (s *Service) UpdateGroup(i UpdateGroupInput, t *jwt.Token) (*GroupBriefOutput, error) {
	uid := t.Claims.(jwt.MapClaims)["uid"].(float64)
	if err := pkg.ValidateStruct(i); err != nil {
		log.Printf("⚠️ Error updating group: %s", err)
		return nil, err
	}
	groups, err := s.repo.GetGroupsByUserID(int(uid))
	if err != nil {
		return nil, err
	}
	userIsInGroup := slices.ContainsFunc(groups, func(g Group) bool {
		return g.ID == i.ID
	})

	if !userIsInGroup {
		err := &e.ErrForbidden{}
		log.Printf("⚠️ Error updating group: %s", err)
		return nil, err
	}

	g, err := s.repo.UpdateGroup(i.ID, i.Name)
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

func (s *Service) UpdateUserGroup(i JoinGroupInput, t *jwt.Token) (*UserGroupOutput, error) {
	uid := t.Claims.(jwt.MapClaims)["uid"].(float64)
	if err := pkg.ValidateStruct(i); err != nil {
		log.Printf("⚠️ Error joining group: %s", err)
		return nil, err
	}
	var ug *UserGroup
	var err error
	switch i.Status {
	case "Accepted":
		ug, err = s.repo.AcceptGroupRequest(i.GroupID, int(uid))
	case "Rejected":
		ug, err = s.repo.RejectGroupRequest(i.GroupID, int(uid))
	case "Invited":
		ug, err = s.repo.InviteGroupRequest(i.GroupID, i.UserID)
	case "Left":
		ug, err = s.repo.LeaveGroup(i.GroupID, int(uid))
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
	}
	if err = pkg.ValidateStruct(out); err != nil {
		log.Printf("⚠️ Error joining group: %s", err)
		return nil, err
	}
	return out, nil
}
