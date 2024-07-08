package auth

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/pkg"
)

type Service struct {
	repo        AuthRepository
	persistence persistence.Persistence
}

func NewService() Service {
	p := persistence.NewConnection()
	return Service{repo: *NewAuthRepository(p.Session()), persistence: *p}
}

func (s *Service) Close() {
	s.persistence.Close()
}

func (s *Service) Login(i LoginInput, rm pkg.RequestMeta) (User, error) {
	err := i.Validate()
	if err != nil {
		log.Printf("⚠️ Error logging in user: %s", err)
		return User{}, err
	}

	uid, err := s.repo.getUserID(i.Email)
	if err != nil {
		log.Printf("⚠️ Error logging in user: %s", err)
		err := &ErrInvalidInput{}
		return User{}, err
	}
	isLocked, err := s.repo.isAccountLocked(uid)
	if err != nil {
		return User{}, fmt.Errorf("unexpected error: %w", err)
	}
	if isLocked {
		err := &ErrAccountLocked{}
		log.Printf("⚠️ Error logging in user: %s", err)
		return User{}, err
	}
	salt, algorithm, err := s.repo.getSaltAndAlgorithmForUser(uid)
	if err != nil {
		return User{}, fmt.Errorf("unexpected error: %w", err)
	}
	hashedPassword := pkg.HashPassword(i.Password, salt, algorithm.Name())
	user, err := s.repo.getInternalLoginUser(i.Email, hashedPassword)
	if err != nil {
		log.Printf("⚠️ Error logging in user: %s", err)
		err := &ErrInvalidInput{}
		s.repo.increaseLoginAttempts(uid)
		return User{}, err
	}
	s.logUserSession(uid, rm)
	return user, nil
}

func (s *Service) Register(i RegisterInput, rm pkg.RequestMeta) (User, error) {
	err := i.Validate()
	if err != nil {
		log.Printf("⚠️ Error registering user: %s", err)
		return User{}, err
	}
	_, err = s.repo.getUserID(i.Email)
	if err == sql.ErrNoRows {
		salt := pkg.GenerateSalt()
		hashedPassword := pkg.HashPassword(i.Password, salt, SHA256.Name())
		user, err := s.repo.createUser(i.Username, i.Email)
		if err != nil {
			log.Printf("⚠️ Error creating user: %s", err)
			return User{}, err
		}
		il, err := s.repo.createInternalLogin(user.ID, hashedPassword, salt, uint16(SHA256))
		if err != nil {
			log.Printf("⚠️ Error creating internal user: %s", err)
			return User{}, err
		}
		user.internalLoginData = il
		log.Printf("✅ User %v created successfully", user.ID)
		s.logUserSession(user.ID, rm)
		return user, nil
	}
	if err != nil {
		return User{}, err
	}
	return User{}, &ErrUserAlreadyExists{}
}

func (s *Service) Me(token *jwt.Token) (User, error) {
	userHash := token.Claims.(jwt.MapClaims)["sub"].(string)
	user, err := s.repo.getUserByHash(userHash)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Service) logUserSession(uid int, rm pkg.RequestMeta) error {
	return s.repo.logUserSession(uid, rm.Ip, rm.Agent)
}
