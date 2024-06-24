package main

import (
	"flag"
	"fmt"

	"github.com/prulloac/fineasy/internal/persistence"
)

func main() {
	p := persistence.NewConnection()
	flag.Parse()
	defer p.Close()
	if flag.Arg(0) == "down" {
		p.DropSchema()
		fmt.Println("Database schema has been dropped.")
		return
	} else if flag.Arg(0) == "reset" {
		p.DropSchema()
		p.VerifySchema()
	} else if flag.Arg(0) == "up" {
		p.VerifySchema()
	}
	fmt.Println("Database schema is up to date.")
}
