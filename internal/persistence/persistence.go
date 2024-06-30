package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	godotenv "github.com/joho/godotenv"

	_ "github.com/lib/pq"
	a "github.com/prulloac/fineasy/internal/auth"
	c "github.com/prulloac/fineasy/internal/currencies"
	n "github.com/prulloac/fineasy/internal/notifications"
	t "github.com/prulloac/fineasy/internal/transactions"
	p "github.com/prulloac/fineasy/internal/user_preferences"
)

type Persistence struct {
	db                        *sql.DB
	authRepository            *a.AuthRepository
	currencyRepository        *c.CurrencyRepository
	transactionsRepository    *t.TransactionsRepository
	notificationsRepository   *n.NotificationsRepository
	userPreferencesRepository *p.UserPreferencesRepository
}

func (p *Persistence) GetAuthRepository() *a.AuthRepository {
	return p.authRepository
}

func (p *Persistence) GetCurrencyRepository() *c.CurrencyRepository {
	return p.currencyRepository
}

func (p *Persistence) GetTransactionsRepository() *t.TransactionsRepository {
	return p.transactionsRepository
}

func (p *Persistence) GetNotificationsRepository() *n.NotificationsRepository {
	return p.notificationsRepository
}

func (p *Persistence) GetUserPreferencesRepository() *p.UserPreferencesRepository {
	return p.userPreferencesRepository
}

func NewConnection() *Persistence {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Database Successfully connected!")

	instance := &Persistence{}
	instance.db = db
	instance.authRepository = a.NewAuthRepository(db)
	instance.currencyRepository = c.NewCurrencyRepository(db)
	instance.transactionsRepository = t.NewTransactionsRepository(db)
	instance.notificationsRepository = n.NewNotificationsRepository(db)
	instance.userPreferencesRepository = p.NewUserPreferencesRepository(db)
	return instance
}

func (p *Persistence) Close() {
	p.db.Close()
	fmt.Println("Database Successfully disconnected!")
}
