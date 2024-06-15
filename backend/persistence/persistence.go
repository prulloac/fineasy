package persistence

import (
	"database/sql"
	"fmt"
	"os"

	godotenv "github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type Persistence struct {
	db *sql.DB
}

func Connect() *Persistence {
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

	return &Persistence{db}
}

func (p *Persistence) Close() {
	p.db.Close()
	fmt.Println("Database Successfully disconnected!")
}

func (p *Persistence) VerifySchema() {
	p.CreateCategoriesTable()
	p.CreateCurrenciesTable()
}
