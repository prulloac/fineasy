package main

import (
	"github.com/prulloac/fineasy/persistence"
	"github.com/prulloac/fineasy/routes"
)

func main() {
	db := persistence.Connection()
	defer db.Close()
	persistence.VerifySchema(db)
	routes.Run()
}
