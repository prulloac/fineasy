package main

import (
	"flag"

	"github.com/prulloac/fineasy/internal/auth"
	"github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/internal/preferences"
	"github.com/prulloac/fineasy/internal/social"
	"github.com/prulloac/fineasy/internal/transactions"
	"github.com/prulloac/fineasy/pkg/logging"
)

type schemaOperation interface {
	CreateTables() error
	DropTables() error
}

var operations = []schemaOperation{}

func main() {
	p := persistence.NewPersistence()
	operations = []schemaOperation{auth.NewRepository(p),
		social.NewRepository(p),
		transactions.NewRepository(p),
		preferences.NewRepository(p),
	}
	flag.Parse()
	defer p.Close()
	if flag.Arg(0) == "down" {
		schemaDown(p)
		logging.Println("Database schema has been dropped.")
		return
	} else if flag.Arg(0) == "reset" {
		schemaDown(p)
		schemaUp(p)
	} else if flag.Arg(0) == "up" {
		schemaUp(p)
	}
	logging.Println("Database schema is up to date.")
}

func schemaUp(p *persistence.Persistence) {
	for _, operation := range operations {
		err := operation.CreateTables()
		if err != nil {
			panic(err)
		}
	}
}

func schemaDown(p *persistence.Persistence) {
	for _, operation := range operations {
		err := operation.DropTables()
		if err != nil {
			panic(err)
		}
	}
}
