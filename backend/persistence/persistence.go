package persistence

import (
	"database/sql"
	"fmt"
	"os"

	godotenv "github.com/joho/godotenv"

	_ "github.com/lib/pq"
	. "github.com/prulloac/fineasy/persistence/repositories"
)

type Persistence struct {
	DB                     *sql.DB
	userRepository         *UserRepository
	categoriesRepository   *CategoriesRepository
	currencyRepository     *CurrencyRepository
	exchangeRateRepository *ExchangeRateRepository
	groupRepository        *GroupRepository
}

var instance *Persistence

func (p *Persistence) GetUserRepository() *UserRepository {
	return p.userRepository
}

func (p *Persistence) GetCategoriesRepository() *CategoriesRepository {
	return p.categoriesRepository
}

func (p *Persistence) GetCurrencyRepository() *CurrencyRepository {
	return p.currencyRepository
}

func (p *Persistence) GetExchangeRateRepository() *ExchangeRateRepository {
	return p.exchangeRateRepository
}

func (p *Persistence) GetGroupRepository() *GroupRepository {
	return p.groupRepository
}

func NewPersistence(db *sql.DB) *Persistence {
	return &Persistence{db,
		&UserRepository{db},
		&CategoriesRepository{db},
		&CurrencyRepository{db},
		&ExchangeRateRepository{db},
		&GroupRepository{db},
	}
}

func Connect() *Persistence {
	if instance != nil {
		return instance
	}
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
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

	instance = NewPersistence(db)
	return instance
}

func (p *Persistence) Close() {
	p.DB.Close()
	fmt.Println("Database Successfully disconnected!")
}

func (p *Persistence) VerifySchema() {
	fmt.Println("Verifying schema...")
	p.GetUserRepository().CreateUsersTable()
	p.GetCurrencyRepository().CreateCurrenciesTable()
	p.GetGroupRepository().CreateGroupsTable()
	p.GetExchangeRateRepository().CreateExchangeRatesTable()
	p.GetCategoriesRepository().CreateCategoriesTable()
	fmt.Println("Schema verified!")
}
