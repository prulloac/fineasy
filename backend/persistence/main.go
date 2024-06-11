package persistence

import (
	"database/sql"
	"fmt"
	"os"

	godotenv "github.com/joho/godotenv"

	_ "github.com/lib/pq"
	"github.com/prulloac/fineasy/persistence/categories"
)

func Connection() *sql.DB {
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

	return db
}

func VerifySchema(db *sql.DB) {
	categories.CreateCategoriesTable(db)
}
