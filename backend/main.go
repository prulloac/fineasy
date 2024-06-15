package main

import (
	"fmt"

	"github.com/prulloac/fineasy/persistence"
	"github.com/prulloac/fineasy/routes"
)

func main() {
	p := persistence.Connect()
	defer p.Close()
	p.VerifySchema()
	fmt.Println("Server is running...")
	routes.Run()
}
