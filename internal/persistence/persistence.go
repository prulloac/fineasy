package persistence

import (
	"database/sql"
	"log"
	"os"

	godotenv "github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
)

type Persistence struct {
	gdb *gorm.DB
}

func NewPersistence() *Persistence {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sql, err := db.DB()
	if err != nil {
		panic(err)
	}
	err = sql.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Database Successfully connected!")

	instance := &Persistence{}
	instance.gdb = db
	return instance
}

func (p *Persistence) Close() {
	sql, err := p.gdb.DB()
	if err != nil {
		panic(err)
	}
	err = sql.Close()
	log.Println("Database Successfully disconnected!")
}

func (p *Persistence) ORM() *gorm.DB {
	return p.gdb
}

func (p *Persistence) SQL() *sql.DB {
	sql, err := p.gdb.DB()
	if err != nil {
		panic(err)
	}
	return sql
}
