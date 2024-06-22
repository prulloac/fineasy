package main

import (
	"fmt"

	server "github.com/prulloac/fineasy/internal/routes"
)

func main() {
	fmt.Println("Server is running...")
	server.Run()
}
