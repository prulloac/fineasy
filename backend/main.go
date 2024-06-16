package main

import (
	"fmt"

	"github.com/prulloac/fineasy/api"
	"github.com/prulloac/fineasy/persistence"
)

func main() {
	p := persistence.Connect()
	defer p.Close()
	p.VerifySchema()
	fmt.Println("Server is running...")
	api.Run()
}
