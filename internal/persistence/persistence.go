package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	godotenv "github.com/joho/godotenv"

	_ "github.com/lib/pq"
	c "github.com/prulloac/fineasy/internal/currencies"
	r "github.com/prulloac/fineasy/internal/persistence/repositories"
)

type Persistence struct {
	db                     *sql.DB
	userRepository         *r.UserRepository
	categoriesRepository   *r.CategoriesRepository
	currencyRepository     *c.CurrencyRepository
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

func (p *Persistence) GetCurrencyRepository() *c.CurrencyRepository {
	return p.currencyRepository
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

func NewConnection() *Persistence {
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
	instance.currencyRepository = c.NewCurrencyRepository(db)
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
	p.createCurrenciesTables()
	p.GetUserRepository().CreateTable()
	p.GetGroupRepository().CreateTable()
	p.GetUserGroupsRepository().CreateTable()
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
	p.GetUserGroupsRepository().DropTable()
	p.GetGroupRepository().DropTable()
	p.GetUserRepository().DropTable()
	p.dropCurrenciesTables()
	fmt.Println("Schema dropped!")
}

func (e *Persistence) createCurrenciesTables() {
	data, _ := os.ReadFile("internal/currencies/schema/currencies_up.sql")

	if data == nil {
		panic("Error reading currencies schema file!")
	}

	_, err := e.db.Exec(string(data))
	if err != nil {
		fmt.Println("Error creating currencies schema!")
		panic(err)
	}
	fmt.Println("Currencies schema up!")
}

func (e *Persistence) dropCurrenciesTables() {
	data, _ := os.ReadFile("internal/currencies/schema/currencies_down.sql")

	if data == nil {
		panic("Error reading currencies schema file!")
	}

	_, err := e.db.Exec(string(data))
	if err != nil {
		fmt.Println("Error dropping currencies schema!")
		panic(err)
	}
	fmt.Println("Currencies schema dropped!")
}
