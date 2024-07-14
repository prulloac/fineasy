package main

import (
	"flag"
	"log"

	"github.com/prulloac/fineasy/internal/auth"
	"github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/internal/social"
)

func main() {
	p := persistence.NewPersistence()
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
	CreateTables() error
	DropTables() error
}

func schemaUp(p *persistence.Persistence) {
	var operations = []schemaOperation{
		auth.NewAuthRepository(p),
		social.NewSocialRepository(p),
	}

	for _, operation := range operations {
		err := operation.CreateTables()
		if err != nil {
			panic(err)
		}
	}
}

func schemaDown(p *persistence.Persistence) {
	var operations = []schemaOperation{
		auth.NewAuthRepository(p),
		social.NewSocialRepository(p),
	}

	for _, operation := range operations {
		err := operation.DropTables()
		if err != nil {
			panic(err)
		}
	}
}
