package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	godotenv "github.com/joho/godotenv"

	_ "github.com/lib/pq"
	r "github.com/prulloac/fineasy/internal/persistence/repositories"
)

type Persistence struct {
	db                     *sql.DB
	userRepository         *r.UserRepository
	categoriesRepository   *r.CategoriesRepository
	currencyRepository     *r.CurrencyRepository
	exchangeRateRepository *r.ExchangeRateRepository
	groupRepository        *r.GroupRepository
}

func (p *Persistence) GetUserRepository() *r.UserRepository {
	return p.userRepository
}

func (p *Persistence) GetCategoriesRepository() *r.CategoriesRepository {
	return p.categoriesRepository
}

func (p *Persistence) GetCurrencyRepository() *r.CurrencyRepository {
	return p.currencyRepository
}

func (p *Persistence) GetExchangeRateRepository() *r.ExchangeRateRepository {
	return p.exchangeRateRepository
}

func (p *Persistence) GetGroupRepository() *r.GroupRepository {
	return p.groupRepository
}

func Connect() *Persistence {
	err := godotenv.Load()

	if err != nil {
		log.Panicln("Error loading .env file")
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
	instance.userRepository = r.NewUserRepository(db)
	instance.categoriesRepository = r.NewCategoriesRepository(db)
	instance.currencyRepository = r.NewCurrencyRepository(db)
	instance.exchangeRateRepository = r.NewExchangeRateRepository(db)
	instance.groupRepository = r.NewGroupRepository(db)
	return instance
}

func (p *Persistence) Close() {
	p.db.Close()
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

func (p *Persistence) DropSchema() {
	fmt.Println("Dropping schema...")
	p.GetCategoriesRepository().DropCategoriesTable()
	p.GetExchangeRateRepository().DropExchangeRatesTable()
	p.GetGroupRepository().DropGroupsTable()
	p.GetCurrencyRepository().DropCurrenciesTable()
	p.GetUserRepository().DropUsersTable()
	fmt.Println("Schema dropped!")
}
