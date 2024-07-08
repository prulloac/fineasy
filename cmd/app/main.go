package main

import (
	"log"

	server "github.com/prulloac/fineasy/internal/routes"
)

func main() {
	log.Println("Server is running...")
	server.Run()
}
