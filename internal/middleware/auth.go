package middleware

import (
	"database/sql"
	"fmt"

	auth "github.com/prulloac/fineasy/internal/db/repositories/auth"
	"github.com/prulloac/fineasy/internal/errors"
	"github.com/prulloac/fineasy/internal/rest"
	"github.com/prulloac/fineasy/internal/rest/dto"
	"github.com/prulloac/fineasy/pkg"
	"github.com/prulloac/fineasy/pkg/logging"
)

type AuthService struct {
	repository       *auth.AuthRepository
	logger           *logging.Logger
	newUserCallbacks []NewUserCallback
}

func NewAuthService(repository *auth.AuthRepository) *AuthService {
	instance := &AuthService{}
	instance.repository = repository
	instance.logger = logging.NewLoggerWithPrefix("auth-service")
	instance.newUserCallbacks = []NewUserCallback{}
	return instance
}

func (s *AuthService) Close() {
	s.repository.Close()
}

type NewUserCallback func(*auth.User)

func (s *AuthService) Register(mail, pwd string, rm rest.RequestMeta) (*dto.UserRegistrationOutput, error) {
	_, err := s.repository.GetUserIDByEmail(mail)
	if err == sql.ErrNoRows {
		salt := pkg.GenerateSalt()
		hashedPassword := pkg.HashPassword(pwd, salt, pkg.SHA256)
		user, err := s.repository.CreateUser(mail)
		if err != nil {
			s.logger.Printf("⚠️ Error creating user: %s", err)
			return nil, err
		}
		il, err := s.repository.CreateInternalLogin(user.ID, hashedPassword, salt, pkg.SHA256)
		if err != nil {
			s.logger.Printf("⚠️ Error creating internal user: %s", err)
			return nil, err
		}
		user.InternalLoginData = *il
		s.logger.Printf("✅ User %v created successfully", user.ID)
		s.logUserSession(user.ID, rm)

		for _, f := range s.newUserCallbacks {
			f(user)
		}

		return &dto.UserRegistrationOutput{
			ID:    user.ID,
			Hash:  user.Hash,
			Email: user.Email,
		}, nil
	}
	if err != nil {
		s.logger.Printf("⚠️ Error creating user: %s", err)
		return nil, err
	}
	return nil, &errors.ErrUserAlreadyExists{}
}

func (s *AuthService) Login(mail, pwd string, rm rest.RequestMeta) (*dto.UserLoginOutput, *auth.User, error) {
	uid, err := s.repository.GetUserIDByEmail(mail)
	if err != nil {
		s.logger.Printf("⚠️ Error logging in user: %s", err)
		err := &errors.ErrInvalidInput{}
		return nil, nil, err
	}
	isLocked, err := s.repository.IsAccountLocked(uid)
	if err != nil {
		return nil, nil, fmt.Errorf("unexpected error: %w", err)
	}
	if isLocked {
		err := &errors.ErrAccountLocked{}
		s.logger.Printf("⚠️ Error logging in user: %s", err)
		return nil, nil, err
	}
	internalUser, err := s.repository.GetInternalLoginByUserID(uid)
	if err != nil {
		return nil, nil, fmt.Errorf("unexpected error: %w", err)
	}
	hashedPassword := pkg.HashPassword(pwd, internalUser.PasswordSalt, internalUser.Algorithm)
	user, err := s.repository.GetUserByEmailAndPassword(mail, hashedPassword)
	if err != nil {
		s.logger.Printf("⚠️ Error logging in user: %s", err)
		err := &errors.ErrInvalidInput{}
		s.repository.IncreaseLoginAttempts(uid)
		return nil, nil, err
	}
	sesh, err := s.logUserSession(uid, rm)
	if err != nil {
		return nil, nil, fmt.Errorf("unexpected error: %w", err)
	}
	return &dto.UserLoginOutput{
		SessionID: sesh.SessionToken,
	}, user, nil
}

func (s *AuthService) logUserSession(uid uint, rm rest.RequestMeta) (*auth.UserSession, error) {
	return s.repository.LogUserSession(uid, rm.Ip, rm.Agent)
}

func (s *AuthService) NewUserCallbacks(f ...NewUserCallback) *AuthService {
	s.newUserCallbacks = append(s.newUserCallbacks, f...)
	return s
}
