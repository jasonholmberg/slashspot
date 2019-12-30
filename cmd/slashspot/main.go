package main

import (
	"log"

	"github.com/joho/godotenv"
	spot "github.com/jasonholmberg/slashspot/internal"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	spot.Run()
}
