package main

import (
	server "github.com/prulloac/fineasy/internal/routes"
	"github.com/prulloac/fineasy/pkg/logging"
)

func main() {
	logging.Println("Server is running...")
	server.Server().Run(":8080")
}
