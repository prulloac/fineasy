package auth

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	e "github.com/prulloac/fineasy/internal/errors"
	p "github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/pkg"
)

type NewUserCallback func(User)

type Service struct {
	repo             *AuthRepository
	newUserCallbacks []NewUserCallback
	logger           *log.Logger
}

func NewService(per *p.Persistence) *Service {
	instance := &Service{}
	instance.repo = NewAuthRepository(per)
	instance.logger = log.New(os.Stdout, "[AuthService] ", log.LUTC)
	return instance
}

func (s *Service) Close() {
	s.repo.Close()
}

func (s *Service) Register(uname, mail, pwd string, rm pkg.RequestMeta) (*UserRegistrationOutput, error) {
	_, err := s.repo.getUserIDByEmail(mail)
	if err == sql.ErrNoRows {
		salt := pkg.GenerateSalt()
		hashedPassword := pkg.HashPassword(pwd, salt, pkg.SHA256)
		user, err := s.repo.createUser(uname, mail)
		if err != nil {
			s.logger.Printf("⚠️ Error creating user: %s", err)
			return nil, err
		}
		il, err := s.repo.createInternalLogin(user.ID, hashedPassword, salt, pkg.SHA256)
		if err != nil {
			s.logger.Printf("⚠️ Error creating internal user: %s", err)
			return nil, err
		}
		user.InternalLoginData = il
		s.logger.Printf("✅ User %v created successfully", user.ID)
		s.logUserSession(user.ID, rm)

		for _, f := range s.newUserCallbacks {
			f(user)
		}

		return &UserRegistrationOutput{
			ID:       user.ID,
			Hash:     user.Hash,
			Username: user.Username,
			Email:    user.Email,
		}, nil
	}
	if err != nil {
		s.logger.Printf("⚠️ Error creating user: %s", err)
		return nil, err
	}
	return nil, &e.ErrUserAlreadyExists{}
}

func (s *Service) Login(mail, pwd string, rm pkg.RequestMeta) (*UserLoginOutput, User, error) {
	uid, err := s.repo.getUserIDByEmail(mail)
	if err != nil {
		s.logger.Printf("⚠️ Error logging in user: %s", err)
		err := &e.ErrInvalidInput{}
		return nil, User{}, err
	}
	isLocked, err := s.repo.isAccountLocked(uid)
	if err != nil {
		return nil, User{}, fmt.Errorf("unexpected error: %w", err)
	}
	if isLocked {
		err := &e.ErrAccountLocked{}
		s.logger.Printf("⚠️ Error logging in user: %s", err)
		return nil, User{}, err
	}
	salt, algorithm, err := s.repo.getSaltAndAlgorithmByUserID(uid)
	if err != nil {
		return nil, User{}, fmt.Errorf("unexpected error: %w", err)
	}
	hashedPassword := pkg.HashPassword(pwd, salt, algorithm)
	user, err := s.repo.getInternalLoginUserByEmailAndPassword(mail, hashedPassword)
	if err != nil {
		s.logger.Printf("⚠️ Error logging in user: %s", err)
		err := &e.ErrInvalidInput{}
		s.repo.increaseLoginAttempts(uid)
		return nil, User{}, err
	}
	sesh, err := s.logUserSession(uid, rm)
	if err != nil {
		return nil, User{}, fmt.Errorf("unexpected error: %w", err)
	}
	return &UserLoginOutput{
		SessionID: sesh.SessionToken,
	}, user, nil
}

func (s *Service) Me(uid uint) (UserOutput, error) {
	user, err := s.repo.getUserByID(uid)
	return UserOutput{
		ID:          user.ID,
		Hash:        user.Hash,
		Username:    user.Username,
		Email:       user.Email,
		Preferences: map[string]interface{}{},
	}, err
}

func (s *Service) GetUserFromToken(token *jwt.Token) (UserOutput, error) {
	uid, ok := token.Claims.(jwt.MapClaims)["uid"].(float64)
	if !ok {
		return UserOutput{}, &e.ErrInvalidInput{}
	}
	return s.Me(uint(uid))
}

func (s *Service) logUserSession(uid uint, rm pkg.RequestMeta) (UserSession, error) {
	return s.repo.logUserSession(uid, rm.Ip, rm.Agent)
}

func (s *Service) NewUserCallbacks(f ...NewUserCallback) *Service {
	s.newUserCallbacks = append(s.newUserCallbacks, f...)
	return s
}
