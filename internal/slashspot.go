package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jasonholmberg/slashspot/internal/handlers"
	"github.com/jasonholmberg/slashspot/internal/store"
)

// Run - Run spot bot, run
func Run() {
	store.Open()
	http.HandleFunc("/command", handlers.SlashCommandHandler)
	port := os.Getenv("SPOT_SERVER_PORT")
	log.Println("Spot's listening on", port)
	http.ListenAndServe(fmt.Sprint(":", port), nil)
}
