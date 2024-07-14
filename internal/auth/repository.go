package auth

import (
	"log"

	p "github.com/prulloac/fineasy/internal/persistence"
	"gorm.io/gorm"
)

type AuthRepository struct {
	Persistence *p.Persistence
}

func NewAuthRepository(persistence *p.Persistence) *AuthRepository {
	return &AuthRepository{persistence}
}

func (a *AuthRepository) Close() {
	a.Persistence.Close()
}

func (a *AuthRepository) CreateTables() error {
	err := a.Persistence.ORM().Migrator().CreateTable(&User{}, &InternalLogin{}, &UserSession{})
	if err != nil {
		log.Println("Error migrating tables")
		return err
	}
	return nil
}

func (a *AuthRepository) DropTables() error {
	err := a.Persistence.ORM().Migrator().DropTable(&User{}, &InternalLogin{}, &UserSession{})
	if err != nil {
		log.Println("Error dropping tables")
		return err
	}
	return nil
}

func (a *AuthRepository) getUserIDByEmail(email string) (uint, error) {
	var uid uint
	r := a.Persistence.ORM().Model(&User{}).Select("id").Where("email = ?", email).First(&uid)
	err := r.Error
	if err != nil {
		log.Printf("Error getting user by email: %s for email: %s", err, email)
		return 0, err
	}
	return uid, nil
}

func (a *AuthRepository) getSaltAndAlgorithmByUserID(uid uint) (string, Algorithm, error) {
	// salt, algorithm,
	var sa struct {
		PasswordSalt string
		Algorithm    Algorithm
	}
	err := a.Persistence.ORM().Model(&InternalLogin{}).Select("password_salt", "algorithm").Where("user_id = ?", uid).Scan(&sa).Error
	return sa.PasswordSalt, sa.Algorithm, err
}

func (a *AuthRepository) getInternalLoginUserByEmailAndPassword(email string, hashedPassword string) (User, error) {
	var user User
	err := a.Persistence.ORM().Model(&User{}).
		Joins("JOIN internal_logins ON users.id = internal_logins.user_id").
		Where("users.email = ? AND internal_logins.password = ?", email, hashedPassword).
		First(&user).Error
	return user, err
}

func (a *AuthRepository) createUser(username string, email string) (User, error) {
	var user User
	err := a.Persistence.ORM().Create(&User{Username: username, Email: email}).Scan(&user).Error
	return user, err
}

func (a *AuthRepository) createInternalLogin(uid uint, hashedPassword string, salt string, algorithm Algorithm) (InternalLogin, error) {
	var il InternalLogin
	err := a.Persistence.ORM().Create(&InternalLogin{
		UserID:       uid,
		Password:     hashedPassword,
		PasswordSalt: salt,
		Algorithm:    Algorithm(algorithm),
	}).Scan(&il).Error
	return il, err
}

func (a *AuthRepository) increaseLoginAttempts(uid uint) error {
	var attempts int
	err := a.Persistence.ORM().Model(&User{}).Where("id = ?", uid).Update("login_attempts", gorm.Expr("login_attempts + 1")).Pluck("login_attempts", &attempts).Error
	return err
}

func (a *AuthRepository) isAccountLocked(uid uint) (bool, error) {
	var disabled bool
	err := a.Persistence.ORM().Model(&User{}).Where("id = ?", uid).Pluck("disabled", &disabled).Error
	return disabled, err
}

func (a *AuthRepository) logUserSession(uid uint, ip string, userAgent string) error {
	return a.Persistence.ORM().Create(&UserSession{
		UserID:    uid,
		LoginIP:   ip,
		UserAgent: userAgent,
	}).Error
}

func (a *AuthRepository) getUserByID(uid uint) (User, error) {
	var user User
	err := a.Persistence.ORM().First(&user, uid).Error
	return user, err
}
