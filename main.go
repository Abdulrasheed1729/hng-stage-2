package main

import (
	"fmt"
	"hng-stage2/api"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	serverPort := fmt.Sprintf(":%s", port)

	server := api.NewAPIServer(serverPort)

	server.Run()
}
