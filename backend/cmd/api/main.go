package main

import (
	"log"

	"canhbuomxanh/backend/internal/app"
)

func main() {
	server, err := app.New()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	if err := server.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
