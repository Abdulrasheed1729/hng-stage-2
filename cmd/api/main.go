package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	serverPort := fmt.Sprintf(":%s", port)

	server := NewAPIServer(serverPort)

	server.Run()
}
