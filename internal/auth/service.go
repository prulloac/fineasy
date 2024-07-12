package auth

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	e "github.com/prulloac/fineasy/internal/errors"
	"github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/pkg"
)

type Service struct {
	repo        *AuthRepository
	persistence *persistence.Persistence
}

func NewService() *Service {
	instance := &Service{}
	instance.persistence = persistence.NewConnection()
	instance.repo = NewAuthRepository(instance.persistence.Session())
	return instance
}

func (s *Service) Close() {
	s.persistence.Close()
}

func (s *Service) Login(mail, pwd string, rm pkg.RequestMeta) (User, error) {
	uid, err := s.repo.getUserID(mail)
	if err != nil {
		log.Printf("⚠️ Error logging in user: %s", err)
		err := &e.ErrInvalidInput{}
		return User{}, err
	}
	isLocked, err := s.repo.isAccountLocked(uid)
	if err != nil {
		return User{}, fmt.Errorf("unexpected error: %w", err)
	}
	if isLocked {
		err := &e.ErrAccountLocked{}
		log.Printf("⚠️ Error logging in user: %s", err)
		return User{}, err
	}
	salt, algorithm, err := s.repo.getSaltAndAlgorithmForUser(uid)
	if err != nil {
		return User{}, fmt.Errorf("unexpected error: %w", err)
	}
	hashedPassword := pkg.HashPassword(pwd, salt, algorithm.Name())
	user, err := s.repo.getInternalLoginUser(mail, hashedPassword)
	if err != nil {
		log.Printf("⚠️ Error logging in user: %s", err)
		err := &e.ErrInvalidInput{}
		s.repo.increaseLoginAttempts(uid)
		return User{}, err
	}
	s.logUserSession(uid, rm)
	return user, nil
}

func (s *Service) Register(uname, mail, pwd string, rm pkg.RequestMeta) (User, error) {
	_, err := s.repo.getUserID(mail)
	if err == sql.ErrNoRows {
		salt := pkg.GenerateSalt()
		hashedPassword := pkg.HashPassword(pwd, salt, SHA256.Name())
		user, err := s.repo.createUser(uname, mail)
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
	return User{}, &e.ErrUserAlreadyExists{}
}

func (s *Service) Me(uhash string) (User, error) {
	user, err := s.repo.getUserByHash(uhash)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Service) GetUserFromToken(token *jwt.Token) (User, error) {
	uhash := token.Claims.(jwt.MapClaims)["sub"].(string)
	return s.Me(uhash)
}

func (s *Service) logUserSession(uid int, rm pkg.RequestMeta) error {
	return s.repo.logUserSession(uid, rm.Ip, rm.Agent)
}
