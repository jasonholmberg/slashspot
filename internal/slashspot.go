package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jasonholmberg/slashspot/internal/data"
	"github.com/jasonholmberg/slashspot/internal/handlers"
)

// Run - Run spot bot, run
func Run() {
	data.Open()
	http.HandleFunc("/command", handlers.SlashCommandHandler)
	port := os.Getenv("SPOT_SERVER_PORT")
	log.Println("Spot's listening on", port)
	http.ListenAndServe(fmt.Sprint(":", port), nil)
}
