package social

import (
	"log"

	"github.com/golang-jwt/jwt/v5"
	e "github.com/prulloac/fineasy/internal/errors"
	"github.com/prulloac/fineasy/internal/persistence"
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

func (s *Service) AddFriend(i AddFriendInput, t *jwt.Token) (*Friend, error) {
	uid := t.Claims.(jwt.MapClaims)["uid"].(float64)
	if int(uid) != i.UserID {
		err := &e.ErrForbidden{}
		log.Printf("⚠️ Error adding friend: %s", err)
		return nil, err
	}
	return s.repo.AddFriend(i.UserID, i.FriendID)
}
