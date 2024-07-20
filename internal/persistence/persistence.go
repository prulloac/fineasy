package persistence

import (
	"database/sql"
	"os"

	godotenv "github.com/joho/godotenv"
	"github.com/prulloac/fineasy/pkg/logging"

	_ "github.com/lib/pq"
)

type Persistence struct {
	sql           *sql.DB
	logger        *logging.Logger
	debugInstance *Persistence
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

	if os.Getenv("ENV") != "production" {
		db.SetMaxIdleConns(1)
		db.SetMaxOpenConns(1)
		db.SetConnMaxLifetime(0)
	}

	instance := &Persistence{sql: db, logger: logging.NewLogger()}
	instance.logger.Infof("Database Successfully connected!")
	instance.logger.SetDepth(3)
	instance.debugInstance = &Persistence{sql: db, logger: logging.NewLoggerWithLevel(logging.Debug)}
	return instance
}

func (p *Persistence) Close() {
	err := p.sql.Close()
	if err != nil {
		panic(err)
	}
	p.logger.Infof("Database Successfully disconnected!")
}

func (p *Persistence) SQL() *sql.DB {
	return p.sql
}

func (p *Persistence) Debug() *Persistence {
	return p.debugInstance
}

func (p *Persistence) Exec(query string, args ...any) (sql.Result, error) {
	p.logger.Debugf("Executing query: %s", query)
	return p.sql.Exec(query, args...)
}

func (p *Persistence) QueryRow(query string, args ...any) *sql.Row {
	p.logger.Debugf("Querying row: %s", query)
	return p.sql.QueryRow(query, args...)
}

func (p *Persistence) Query(query string, args ...any) (*sql.Rows, error) {
	p.logger.Debugf("Querying: %s", query)
	return p.sql.Query(query, args...)
}
