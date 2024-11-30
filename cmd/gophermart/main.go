package main

import (
	"log"

	"github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart"
)

func main() {
	// Create and initialize the application
	application, err := app.New()
	if err != nil {
		log.Fatalf("failed to initialize application : %s", err.Error())
	}

	// Run the application
	err = application.Run()
	if err != nil {
		log.Fatalf("failed to run application : %s", err.Error())
	}
}
