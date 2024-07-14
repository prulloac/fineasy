package auth

import (
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
	e "github.com/prulloac/fineasy/internal/errors"
	p "github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/pkg"
	"gorm.io/gorm"
)

type Service struct {
	repo            *AuthRepository
	newUserTriggers []func(User)
}

func NewService(per *p.Persistence) *Service {
	instance := &Service{}
	instance.repo = NewAuthRepository(per)
	return instance
}

func (s *Service) Close() {
	s.repo.Close()
}

func (s *Service) Register(uname, mail, pwd string, rm pkg.RequestMeta) (User, error) {
	_, err := s.repo.getUserIDByEmail(mail)
	if err == gorm.ErrRecordNotFound {
		salt := pkg.GenerateSalt()
		hashedPassword := pkg.HashPassword(pwd, salt, SHA256.String())
		user, err := s.repo.createUser(uname, mail)
		if err != nil {
			log.Printf("⚠️ Error creating user: %s", err)
			return User{}, err
		}
		il, err := s.repo.createInternalLogin(user.ID, hashedPassword, salt, SHA256)
		if err != nil {
			log.Printf("⚠️ Error creating internal user: %s", err)
			return User{}, err
		}
		user.InternalLoginData = il
		log.Printf("✅ User %v created successfully", user.ID)
		s.logUserSession(user.ID, rm)

		for _, f := range s.newUserTriggers {
			f(user)
		}

		return user, nil
	}
	if err != nil {
		log.Printf("⚠️ Error creating user: %s", err)
		return User{}, err
	}
	return User{}, &e.ErrUserAlreadyExists{}
}

func (s *Service) Login(mail, pwd string, rm pkg.RequestMeta) (User, error) {
	uid, err := s.repo.getUserIDByEmail(mail)
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
	salt, algorithm, err := s.repo.getSaltAndAlgorithmByUserID(uid)
	if err != nil {
		return User{}, fmt.Errorf("unexpected error: %w", err)
	}
	hashedPassword := pkg.HashPassword(pwd, salt, algorithm.String())
	user, err := s.repo.getInternalLoginUserByEmailAndPassword(mail, hashedPassword)
	if err != nil {
		log.Printf("⚠️ Error logging in user: %s", err)
		err := &e.ErrInvalidInput{}
		s.repo.increaseLoginAttempts(uid)
		return User{}, err
	}
	s.logUserSession(uid, rm)
	return user, nil
}

func (s *Service) Me(uid uint) (User, error) {
	return s.repo.getUserByID(uid)
}

func (s *Service) GetUserFromToken(token *jwt.Token) (User, error) {
	uid, ok := token.Claims.(jwt.MapClaims)["uid"].(float64)
	if !ok {
		return User{}, &e.ErrInvalidInput{}
	}
	return s.Me(uint(uid))
}

func (s *Service) logUserSession(uid uint, rm pkg.RequestMeta) error {
	return s.repo.logUserSession(uid, rm.Ip, rm.Agent)
}

func (s *Service) AddNewUserTrigger(f func(User)) *Service {
	s.newUserTriggers = append(s.newUserTriggers, f)
	return s
}
