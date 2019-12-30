package main

import (
	"log"

	spot "github.com/jasonholmberg/slashspot/internal"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	spot.Run()
}
