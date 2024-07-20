package preferences

import (
	p "github.com/prulloac/fineasy/internal/persistence"
)

type Service struct {
	repo *Repository
}

func NewService(persistence *p.Persistence) *Service {
	instance := &Service{}
	instance.repo = NewRepository(persistence)
	instance.repo.CreateTables()
	return instance
}

func (s *Service) Close() {
	s.repo.Close()
}

func (s *Service) CreateUserData(userID uint) (*UserDataOutput, error) {
	userData, err := s.repo.CreateUserData(userID)
	if err != nil {
		return nil, err
	}
	return &UserDataOutput{
		ID:          userData.ID,
		UserID:      userData.UserID,
		AvatarURL:   userData.AvatarURL,
		DisplayName: userData.DisplayName,
		Currency:    userData.Currency,
	}, nil
}
