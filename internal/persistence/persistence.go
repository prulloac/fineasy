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
	accountsRepository     *r.AccountsRepository
	budgetsRepository      *r.BudgetsRepository
	userGroupsRepository   *r.UserGroupsRepository
	transactionsRepository *r.TransactionsRepository
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

func (p *Persistence) GetAccountsRepository() *r.AccountsRepository {
	return p.accountsRepository
}

func (p *Persistence) GetBudgetsRepository() *r.BudgetsRepository {
	return p.budgetsRepository
}

func (p *Persistence) GetUserGroupsRepository() *r.UserGroupsRepository {
	return p.userGroupsRepository
}

func (p *Persistence) GetTransactionsRepository() *r.TransactionsRepository {
	return p.transactionsRepository
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
	instance.accountsRepository = r.NewAccountsRepository(db)
	instance.budgetsRepository = r.NewBudgetsRepository(db)
	instance.userGroupsRepository = r.NewUserGroupsRepository(db)
	instance.transactionsRepository = r.NewTransactionsRepository(db)
	return instance
}

func (p *Persistence) Close() {
	p.db.Close()
	fmt.Println("Database Successfully disconnected!")
}

func (p *Persistence) VerifySchema() {
	fmt.Println("Verifying schema...")
	p.GetUserRepository().CreateTable()
	p.GetCurrencyRepository().CreateTable()
	p.GetGroupRepository().CreateTable()
	p.GetUserGroupsRepository().CreateTable()
	p.GetExchangeRateRepository().CreateTable()
	p.GetCategoriesRepository().CreateTable()
	p.GetAccountsRepository().CreateTable()
	p.GetBudgetsRepository().CreateTable()
	p.GetTransactionsRepository().CreateTable()
	fmt.Println("Schema verified!")
}

func (p *Persistence) DropSchema() {
	fmt.Println("Dropping schema...")
	p.GetTransactionsRepository().DropTable()
	p.GetBudgetsRepository().DropTable()
	p.GetAccountsRepository().DropTable()
	p.GetCategoriesRepository().DropTable()
	p.GetExchangeRateRepository().DropTable()
	p.GetUserGroupsRepository().DropTable()
	p.GetGroupRepository().DropTable()
	p.GetCurrencyRepository().DropTable()
	p.GetUserRepository().DropTable()
	fmt.Println("Schema dropped!")
}
