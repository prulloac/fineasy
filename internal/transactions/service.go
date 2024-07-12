package transactions

import (
	"log"
	"slices"
	"time"

	e "github.com/prulloac/fineasy/internal/errors"
	"github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/internal/social"
	"github.com/prulloac/fineasy/pkg"
)

type Service struct {
	repo        *TransactionsRepository
	persistence *persistence.Persistence
	social      *social.SocialRepository
}

func NewService() *Service {
	instance := &Service{}
	instance.persistence = persistence.NewConnection()
	instance.repo = NewTransactionsRepository(instance.persistence.Session())
	instance.social = social.NewSocialRepository(instance.persistence.Session())
	return instance
}

func (s *Service) Close() {
	s.persistence.Close()
}

func (s *Service) CreateAccount(name string, groupID int, currency string, uid int) (*CreateAccountOutput, error) {
	ugs, err := s.social.GetUserGroupsByUserID(uid)
	if !slices.ContainsFunc(ugs, func(i social.UserGroup) bool {
		return i.GroupID == groupID
	}) {
		err := &e.ErrForbidden{}
		log.Printf("⚠️ Error creating account: %s", err)
		return nil, err
	}
	a, err := s.repo.CreateAccount(name, groupID, currency, uid)
	if err != nil {
		return nil, err
	}
	out := &CreateAccountOutput{
		ID:        a.ID,
		Name:      a.Name,
		GroupID:   a.GroupID,
		Currency:  a.Currency,
		Balance:   a.Balance,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
	}
	if err = pkg.ValidateStruct(out); err != nil {
		log.Printf("⚠️ Error creating account: %s", err)
		return nil, err
	}
	return out, nil
}
