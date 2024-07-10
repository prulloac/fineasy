package social

import (
	"log"

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
	return s.repo.AddFriend(i.UserID, i.FriendID)
}

func (s *Service) GetFriends(t *jwt.Token) ([]FriendShipOutput, error) {
	uid := t.Claims.(jwt.MapClaims)["uid"].(float64)
	return s.repo.GetFriends(int(uid))
}

func (s *Service) GetFriendRequests(t *jwt.Token) ([]FriendRequestOutput, error) {
	uid := t.Claims.(jwt.MapClaims)["uid"].(float64)
	return s.repo.GetFriendRequests(int(uid))
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
	if i.Status == "Accepted" {
		return s.repo.AcceptFriendRequest(i.UserID, i.FriendID)
	}
	return s.repo.RejectFriendRequest(i.UserID, i.FriendID)
}

func (s *Service) CreateGroup(i CreateGroupInput, t *jwt.Token) (*GroupOutput, error) {
	uid := t.Claims.(jwt.MapClaims)["uid"].(float64)
	if !pkg.Contains(i.Members, int(uid)) {
		err := &e.ErrForbidden{}
		log.Printf("⚠️ Error creating group: %s", err)
		return nil, err
	}
	return s.repo.CreateGroup(i.Name, i.Members, int(uid))
}
