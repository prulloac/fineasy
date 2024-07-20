package transactions

import (
	"slices"
	"strconv"
	"time"

	e "github.com/prulloac/fineasy/internal/errors"
	p "github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/internal/social"
)

type Service struct {
	repo   *Repository
	social *social.Repository
}

func NewService(persistence *p.Persistence) *Service {
	instance := &Service{}
	instance.repo = NewRepository(persistence)
	instance.social = social.NewRepository(persistence)
	instance.repo.CreateTables()
	return instance
}

func (s *Service) Close() {
	s.repo.Close()
	s.social.Close()
}

func (s *Service) CreateAccount(name string, currency string, gid, uid uint) (*CreateAccountOutput, error) {
	ugs, err := s.social.GetUserGroupsByUserID(uid)
	if err != nil {
		return nil, err
	}
	if !slices.ContainsFunc(ugs, func(i social.UserGroup) bool {
		return i.GroupID == gid
	}) {
		err := &e.ErrForbidden{}
		return nil, err
	}
	a, err := s.repo.CreateAccount(name, currency, gid, uid)
	if err != nil {
		return nil, err
	}
	return &CreateAccountOutput{
		ID:        a.ID,
		Name:      a.Name,
		GroupID:   a.GroupID,
		Currency:  a.Currency,
		Balance:   a.Balance,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
	}, nil
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

func (s *Service) GetAccountByID(id, uid uint) (*AccountBriefOutput, error) {
	ok, err := s.repo.UserHasAccessToAccount(uid, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		err := &e.ErrForbidden{}
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

func (s *Service) UpdateAccount(name string, cur string, balance float64, id, uid uint) (*AccountBriefOutput, error) {
	ok, err := s.repo.UserHasAccessToAccount(uid, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		err := &e.ErrForbidden{}
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

func (s *Service) CreateBudget(name, cur string, amount float64, start, end time.Time, aid, uid uint) (*BudgetOutput, error) {
	ok, err := s.repo.UserHasAccessToAccount(uid, aid)
	if err != nil {
		return nil, err
	}
	if !ok {
		err := &e.ErrForbidden{}
		return nil, err
	}
	b, err := s.repo.CreateBudget(name, cur, amount, start, end, aid, uid)
	if err != nil {
		return nil, err
	}
	return &BudgetOutput{
		ID:        b.ID,
		Name:      b.Name,
		AccountID: b.AccountID,
		Currency:  b.Currency,
		Amount:    strconv.FormatFloat(b.Amount, 'f', 2, 64),
		StartDate: b.StartDate.Format(time.RFC3339),
		EndDate:   b.EndDate.Format(time.RFC3339),
	}, nil
}

func (s *Service) GetBudgets(uid uint) ([]BudgetOutput, error) {
	bs, err := s.repo.GetBudgetsByUserID(uid)
	if err != nil {
		return nil, err
	}
	var out []BudgetOutput
	for _, b := range bs {
		out = append(out, BudgetOutput{
			ID:        b.ID,
			Name:      b.Name,
			AccountID: b.AccountID,
			Currency:  b.Currency,
			Amount:    strconv.FormatFloat(b.Amount, 'f', 2, 64),
			StartDate: b.StartDate.Format(time.DateOnly),
			EndDate:   b.EndDate.Format(time.DateOnly),
		})
	}
	return out, nil
}
