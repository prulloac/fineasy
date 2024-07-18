package persistence

import (
	"database/sql"
	"log"
	"os"

	godotenv "github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type Persistence struct {
	db     *sql.DB
	logger *log.Logger
}

func NewPersistence() *Persistence {
	godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	instance := &Persistence{db: db, logger: log.New(os.Stdout, "[Persistence] ", log.LUTC)}

	instance.logger.Println("Database Successfully connected!")
	return instance
}

func (p *Persistence) Close() {
	err := p.db.Close()
	if err != nil {
		panic(err)
	}
	p.logger.Println("Database Successfully disconnected!")
}

func (p *Persistence) SQL() *sql.DB {
	return p.db
}
