package middleware

import (
	"slices"
	"strconv"
	"time"

	core "github.com/prulloac/fineasy/internal/db/repositories/core"
	"github.com/prulloac/fineasy/internal/errors"
	"github.com/prulloac/fineasy/internal/rest/dto"
	"github.com/prulloac/fineasy/pkg/logging"
)

type CoreService struct {
	repository       *core.CoreRepository
	logger           *logging.Logger
	newUserCallbacks []NewUserCallback
}

func NewCoreService(repository *core.CoreRepository) *CoreService {
	instance := &CoreService{}
	instance.repository = repository
	instance.logger = logging.NewLoggerWithPrefix("core-service")
	return instance
}

func (s *CoreService) Close() {
	s.repository.Close()
}

func (s *CoreService) CreateUserData(userID uint) (*dto.UserDataOutput, error) {
	userData, err := s.repository.CreateUserData(userID)
	if err != nil {
		return nil, err
	}
	return &dto.UserDataOutput{
		UserID:      userData.UserID,
		AvatarURL:   userData.AvatarURL,
		DisplayName: userData.DisplayName,
		Currency:    userData.Currency,
	}, nil
}

func (s *CoreService) AddFriendship(fid, uid uint) (*dto.FriendRequestOutput, error) {
	fr, err := s.repository.CreateFriendship(uid, fid)
	if err != nil {
		return nil, err
	}
	return &dto.FriendRequestOutput{
		UserID:   fr.UserID,
		FriendID: fr.FriendID,
		Status:   fr.Status.String(),
	}, nil
}

func (s *CoreService) GetFriendships(uid uint) ([]dto.FriendShipOutput, error) {
	fs, err := s.repository.GetFriendshipsByUserID(uid)
	if err != nil {
		return nil, err
	}
	out := []dto.FriendShipOutput{}
	for _, f := range fs {
		e := dto.FriendShipOutput{
			UserID:       f.UserID,
			FriendID:     f.FriendID,
			RelationType: f.RelationType.String(),
		}
		out = append(out, e)
	}
	return out, nil
}

func (s *CoreService) GetFriendship(fid, uid uint) (*dto.FriendShipOutput, error) {
	f, err := s.repository.GetFriendshipByFriendIDAndUserID(fid, uid)
	if err != nil {
		return nil, err
	}
	return &dto.FriendShipOutput{
		UserID:       f.UserID,
		FriendID:     f.FriendID,
		RelationType: f.RelationType.String(),
	}, nil
}

func (s *CoreService) GetPendingFriendships(uid uint) ([]dto.FriendRequestOutput, error) {
	frs, err := s.repository.GetPendingFriendshipsByUserID(uid)
	if err != nil {
		return nil, err
	}
	out := []dto.FriendRequestOutput{}
	for _, fr := range frs {
		e := dto.FriendRequestOutput{
			UserID:   fr.UserID,
			FriendID: fr.FriendID,
			Status:   fr.Status.String(),
		}
		out = append(out, e)
	}
	return out, nil
}

func (s *CoreService) AcceptFriendship(status string, fid, uid uint) (*dto.FriendRequestOutput, error) {
	var fr *core.Friendship
	var err error
	if status == "Accepted" {
		fr, err = s.repository.AcceptFriendship(uid, fid)
	}
	if err != nil {
		return nil, err
	}
	return &dto.FriendRequestOutput{
		UserID:   fr.UserID,
		FriendID: fr.FriendID,
		Status:   fr.Status.String(),
	}, nil
}

func (s *CoreService) RejectFriendship(fid, uid uint) ([]dto.FriendShipOutput, error) {
	err := s.repository.RejectFriendship(uid, fid)
	if err != nil {
		return nil, err
	}
	return s.GetFriendships(uid)
}

func (s *CoreService) CreateGroup(name string, uid uint) (*dto.GroupBriefOutput, error) {
	g, err := s.repository.CreateGroup(name, uid)
	if err != nil {
		return nil, err
	}
	return &dto.GroupBriefOutput{
		ID:        g.ID,
		Name:      g.Name,
		CreatedBy: g.CreatedBy,
	}, nil
}

func (s *CoreService) GetGroupByID(gid, uid uint) (*dto.GroupFullOutput, error) {
	g, err := s.repository.GetGroupByUserID(gid, uid)
	if err != nil {
		return nil, err
	}
	ms, err := s.repository.GetMembershipsByGroupID(gid)
	if err != nil {
		return nil, err
	}
	out := &dto.GroupFullOutput{
		GroupID:     g.ID,
		Name:        g.Name,
		MemberCount: len(ms),
		CreatedBy:   g.CreatedBy,
		Memberships: []dto.MembershipOutput{},
	}
	for _, m := range ms {
		e := dto.MembershipOutput{
			UserID:   m.UserID,
			Status:   m.Status.String(),
			JoinedAt: m.JoinedAt.String(),
		}
		out.Memberships = append(out.Memberships, e)
	}
	return out, nil
}

func (s *CoreService) GetGroup(id uint) (*dto.GroupBriefOutput, error) {
	g, err := s.repository.GetGroupByID(id)
	if err != nil {
		return nil, err
	}
	return &dto.GroupBriefOutput{
		ID:        g.ID,
		Name:      g.Name,
		CreatedBy: g.CreatedBy,
	}, nil
}

func (s *CoreService) GetUserGroups(uid uint) ([]dto.UserGroupOutput, error) {
	ugs, err := s.repository.GetUserGroupsByUserID(uid)
	if err != nil {
		return nil, err
	}
	out := []dto.UserGroupOutput{}
	for _, ug := range ugs {
		g, err := s.repository.GetGroupByID(ug.GroupID)
		if err != nil {
			return nil, err
		}
		leftAt := ""
		if ug.LeftAt.Valid {
			leftAt = ug.LeftAt.Time.String()
		}
		e := dto.UserGroupOutput{
			UserID:    ug.UserID,
			GroupID:   ug.GroupID,
			Status:    ug.Status.String(),
			Group:     g.Name,
			CreatedBy: g.CreatedBy,
			JoinedAt:  ug.JoinedAt.String(),
			LeftAt:    leftAt,
		}
		out = append(out, e)
	}
	return out, nil
}

func (s *CoreService) UpdateGroup(name string, gid, uid uint) (*dto.GroupBriefOutput, error) {
	groups, err := s.repository.GetGroupsByUserID(uid)
	if err != nil {
		return nil, err
	}
	userIsInGroup := slices.ContainsFunc(groups, func(g core.Group) bool {
		return g.ID == gid
	})

	if !userIsInGroup {
		err := &errors.ErrForbidden{}
		return nil, err
	}

	g, err := s.repository.UpdateGroup(gid, name)
	if err != nil {
		return nil, err
	}
	return &dto.GroupBriefOutput{
		ID:        g.ID,
		Name:      g.Name,
		CreatedBy: g.CreatedBy,
	}, nil
}

func (s *CoreService) UpdateUserGroup(status string, gid, uid uint) (*dto.UserGroupOutput, error) {
	var ug *core.UserGroup
	var err error
	switch status {
	case "Accepted":
		ug, err = s.repository.AcceptGroupRequest(gid, uid)
	case "Rejected":
		ug, err = s.repository.RejectGroupRequest(gid, uid)
	case "Left":
		ug, err = s.repository.LeaveGroup(gid, uid)
	default:
		err := &errors.ErrBadRequest{}
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
	return &dto.UserGroupOutput{
		UserID:    ug.UserID,
		GroupID:   ug.GroupID,
		Status:    ug.Status.String(),
		Group:     g.Name,
		JoinedAt:  ug.JoinedAt.String(),
		LeftAt:    leftAt,
		CreatedBy: g.CreatedBy,
	}, nil
}

func (s *CoreService) InviteUserGroup(gid, uid uint) (*dto.GroupFullOutput, error) {
	ug, err := s.repository.InviteGroupRequest(gid, uid)
	if err != nil {
		return nil, err
	}
	_, err = s.GetGroup(ug.GroupID)
	if err != nil {
		return nil, err
	}
	return s.GetGroupFullData(ug.GroupID)
}

func (s *CoreService) GetGroupFullData(gid uint) (*dto.GroupFullOutput, error) {
	g, err := s.repository.GetGroupByID(gid)
	if err != nil {
		return nil, err
	}
	ms, err := s.repository.GetMembershipsByGroupID(gid)
	if err != nil {
		return nil, err
	}
	out := &dto.GroupFullOutput{
		GroupID:     g.ID,
		Name:        g.Name,
		MemberCount: len(ms),
		CreatedBy:   g.CreatedBy,
		Memberships: []dto.MembershipOutput{},
	}
	for _, m := range ms {
		e := dto.MembershipOutput{
			UserID:   m.UserID,
			Status:   m.Status.String(),
			JoinedAt: m.JoinedAt.String(),
		}
		out.Memberships = append(out.Memberships, e)
	}
	return out, nil
}

func (s *CoreService) CreateAccount(name string, currency string, gid, uid uint) (*dto.CreateAccountOutput, error) {
	a, err := s.repository.CreateAccount(name, currency, gid, uid)
	if err != nil {
		return nil, err
	}
	return &dto.CreateAccountOutput{
		ID:        a.ID,
		Name:      a.Name,
		GroupID:   a.GroupID,
		Currency:  a.Currency,
		Balance:   a.Balance,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *CoreService) GetAccounts(uid uint) ([]dto.AccountBriefOutput, error) {
	as, err := s.repository.GetAccountsByUserID(uid)
	if err != nil {
		return nil, err
	}
	var out []dto.AccountBriefOutput
	for _, a := range as {
		out = append(out, dto.AccountBriefOutput{
			ID:       a.ID,
			Name:     a.Name,
			Currency: a.Currency,
			Balance:  strconv.FormatFloat(a.Balance, 'f', 2, 64),
		})
	}
	return out, nil
}

func (s *CoreService) GetAccountByID(id, uid uint) (*dto.AccountBriefOutput, error) {
	ok, err := s.repository.UserHasAccessToAccount(uid, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		err := &errors.ErrForbidden{}
		return nil, err
	}
	a, err := s.repository.GetAccountByID(id)
	if err != nil {
		return nil, err
	}
	out := &dto.AccountBriefOutput{
		ID:       a.ID,
		Name:     a.Name,
		Currency: a.Currency,
		Balance:  strconv.FormatFloat(a.Balance, 'f', 2, 64),
	}
	return out, nil
}

func (s *CoreService) UpdateAccount(name string, cur string, balance float64, id, uid uint) (*dto.AccountBriefOutput, error) {
	ok, err := s.repository.UserHasAccessToAccount(uid, id)
	if err != nil {
		return nil, err
	}
	if !ok {
		err := &errors.ErrForbidden{}
		return nil, err
	}
	a, err := s.repository.UpdateAccount(id, name, cur, balance)
	if err != nil {
		return nil, err
	}
	out := &dto.AccountBriefOutput{
		ID:       a.ID,
		Name:     a.Name,
		Currency: a.Currency,
		Balance:  strconv.FormatFloat(a.Balance, 'f', 2, 64),
	}
	return out, nil
}

func (s *CoreService) CreateBudget(name, cur string, amount float64, start, end time.Time, aid, uid uint) (*dto.BudgetOutput, error) {
	ok, err := s.repository.UserHasAccessToAccount(uid, aid)
	if err != nil {
		return nil, err
	}
	if !ok {
		err := &errors.ErrForbidden{}
		return nil, err
	}
	b, err := s.repository.CreateBudget(name, cur, amount, start, end, aid, uid)
	if err != nil {
		return nil, err
	}
	return &dto.BudgetOutput{
		ID:        b.ID,
		Name:      b.Name,
		AccountID: b.AccountID,
		Currency:  b.Currency,
		Amount:    strconv.FormatFloat(b.Amount, 'f', 2, 64),
		StartDate: b.StartDate.Format(time.RFC3339),
		EndDate:   b.EndDate.Format(time.RFC3339),
	}, nil
}

func (s *CoreService) GetBudgets(uid uint) ([]dto.BudgetOutput, error) {
	bs, err := s.repository.GetBudgetsByUserID(uid)
	if err != nil {
		return nil, err
	}
	var out []dto.BudgetOutput
	for _, b := range bs {
		out = append(out, dto.BudgetOutput{
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
