package persistence

import (
	"database/sql"
	"log"
	"os"

	godotenv "github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type Persistence struct {
	db *sql.DB
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

	log.Println("Database Successfully connected!")

	instance := &Persistence{}
	instance.db = db
	return instance
}

func (p *Persistence) Close() {
	p.db.Close()
	log.Println("Database Successfully disconnected!")
}

func (p *Persistence) Session() *sql.DB {
	return p.db
}
