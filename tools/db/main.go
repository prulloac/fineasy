package main

import (
	"flag"
	"log"

	"github.com/prulloac/fineasy/internal/auth"
	"github.com/prulloac/fineasy/internal/persistence"
)

func main() {
	p := persistence.NewConnection()
	flag.Parse()
	defer p.Close()
	if flag.Arg(0) == "down" {
		schemaDown(p)
		log.Println("Database schema has been dropped.")
		return
	} else if flag.Arg(0) == "reset" {
		schemaDown(p)
		schemaUp(p)
	} else if flag.Arg(0) == "up" {
		schemaUp(p)
	}
	log.Println("Database schema is up to date.")
}

type schemaOperation interface {
	CreateTable() error
	DropTable() error
}

func schemaUp(p *persistence.Persistence) {
	var operations = []schemaOperation{
		&auth.AuthRepository{DB: p.Session()},
	}

	for _, operation := range operations {
		err := operation.CreateTable()
		if err != nil {
			panic(err)
		}
	}
}

func schemaDown(p *persistence.Persistence) {
	var operations = []schemaOperation{
		&auth.AuthRepository{DB: p.Session()},
	}

	for _, operation := range operations {
		err := operation.DropTable()
		if err != nil {
			panic(err)
		}
	}
}
