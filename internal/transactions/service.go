package transactions

import (
	"log"
	"slices"
	"strconv"
	"time"

	e "github.com/prulloac/fineasy/internal/errors"
	p "github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/internal/social"
	"github.com/prulloac/fineasy/pkg"
)

type Service struct {
	repo   *TransactionsRepository
	social *social.SocialRepository
}

func NewService(persistence *p.Persistence) *Service {
	instance := &Service{}
	instance.repo = NewTransactionsRepository(persistence)
	instance.social = social.NewSocialRepository(persistence)
	return instance
}

func (s *Service) Close() {
	s.repo.Close()
	s.social.Close()
}

func (s *Service) CreateAccount(name string, currency string, gid, uid uint) (*CreateAccountOutput, error) {
	ugs, err := s.social.GetUserGroupsByUserID(uid)
	if !slices.ContainsFunc(ugs, func(i social.UserGroup) bool {
		return i.GroupID == gid
	}) {
		err := &e.ErrForbidden{}
		log.Printf("⚠️ Error creating account: %s", err)
		return nil, err
	}
	a, err := s.repo.CreateAccount(name, currency, gid, uid)
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

func (s *Service) GetAccounts(uid uint) ([]AccountBriefOutput, error) {
	as, err := s.repo.GetAccountsByUserID(uid)
	if err != nil {
		return nil, err
	}
	var out []AccountBriefOutput
	for _, a := range as {
		out = append(out, AccountBriefOutput{
			ID:       a.ID,
			Name:     a.Name,
			Currency: a.Currency,
			Balance:  strconv.FormatFloat(a.Balance, 'f', 2, 64),
		})
	}
	return out, nil
}

func (s *Service) GetAccountByID(id, uid int) (*AccountBriefOutput, error) {
	ok, err := s.repo.UserHasAccessToAccount(uid, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		err := &e.ErrForbidden{}
		log.Printf("⚠️ Error getting account: %s", err)
		return nil, err
	}
	a, err := s.repo.GetAccountByID(id)
	if err != nil {
		return nil, err
	}
	out := &AccountBriefOutput{
		ID:       a.ID,
		Name:     a.Name,
		Currency: a.Currency,
		Balance:  strconv.FormatFloat(a.Balance, 'f', 2, 64),
	}
	return out, nil
}

func (s *Service) UpdateAccount(name string, cur string, balance float64, id, uid int) (*AccountBriefOutput, error) {
	ok, err := s.repo.UserHasAccessToAccount(uid, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		err := &e.ErrForbidden{}
		log.Printf("⚠️ Error updating account: %s", err)
		return nil, err
	}
	a, err := s.repo.UpdateAccount(id, name, cur, balance)
	if err != nil {
		return nil, err
	}
	out := &AccountBriefOutput{
		ID:       a.ID,
		Name:     a.Name,
		Currency: a.Currency,
		Balance:  strconv.FormatFloat(a.Balance, 'f', 2, 64),
	}
	return out, nil
}
