package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/nlopes/slack"
	"github.com/jasonholmberg/slashspot/internal/store"
	"github.com/jasonholmberg/slashspot/internal/util"
)

const (
	// HelpText - the standard help text
	HelpText = `
	*Slash-Spot Help*:
	With Slash-Spot you can find, reserve and register parking spots.

	*/spot find or open* - will deliver a list of spots available today
	*/spot claim or take or reserve <spot-id>* - will attempt claim/take/reserve the requested spot
	*/spot reg or register or set <spot-id> [date]* - will make a spot available for use for the day. 
		If a data is given, the spot will be made available for that date. That date must be in the future.
	*/spot drop <spot-id>* - will attempt to drop a spot registration as long as your are the registering user
	*/spot drop all* - will attempt to drop all spots you have registered.
	`
	// IDKBlank - I don't know message template
	IDKBlank = "I don't know what you mean, use `/spot help` for some...help."

	// IDKTemplate - I don't know message template
	IDKTemplate = "I don't know what '%s' means, use `/spot help` for some...help."

	// OpenSpotsTemplate - Open spots template
	OpenSpotsTemplate = "The follow spots are available today: %v"

	// NoSpotsAvailable - No spots
	NoSpotsAvailable = "There are currently no available registered spots."

	// SpotClaimedTemplate - Spot claimed template
	SpotClaimedTemplate = "You have claimed spot: %v"

	// SpotClaimErrorTemplate - Claim error template
	SpotClaimErrorTemplate = "The spot %s is not available today or has not been registered as available"

	// SpotRegisteredTemplate - Spot registered template
	SpotRegisteredTemplate = "You have registered spot %s. Thank you for sharing"

	// SpotDupeRegistrationErrorTemplate - Spot registration error
	SpotDupeRegistrationErrorTemplate = "The Spot %v has already been register by %v"

	// SpotDateFormatRegistrationErrorTemplate - Spot registration date format error
	SpotDateFormatRegistrationErrorTemplate = "The date provided: %s is invalid, please use format YYYY-MM-DD"

	// SpotPastDateRegistrationErrorTemplate - Spot past date error
	SpotPastDateRegistrationErrorTemplate = "The date provided: %s, is in the past,"

	// SpotDropRegTemplate - Spot drop registration template
	SpotDropRegTemplate = "Registration for spot %s has been dropped"

	// SpotDropRegErrorTemplate - Error respose template for drop registration error
	SpotDropRegErrorTemplate = "Unable to drop registration %v. The registration has been claimed or you did not create this registration."
)

// SlashCommandHandler - the root handler for spot.  Capture the incoming command from slack and delegates it off to other internal handlers.
func SlashCommandHandler(w http.ResponseWriter, r *http.Request) {
	verifier, err := slack.NewSecretsVerifier(r.Header, os.Getenv("SPOT_SLACK_SIGNING_SECRET"))
	if err != nil {
		log.Print("ERROR - slashspot may not be configured correctly, check you set up: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/spot":
		spotCommandHandler(&s, w)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func spotCommandHandler(cmd *slack.SlashCommand, w http.ResponseWriter) {
	params := strings.Split(cmd.Text, " ")
	log.Printf("Spot command received %v", params)
	var response string
	switch action := params[0]; action {
	case "", " ":
		response = handleBlank()
	case "help":
		response = handleHelp()
	case "find", "open":
		response = handleFind(params)
	case "reg", "register", "set":
		response = handleRegister(cmd, params)
	case "claim", "take", "reserve":
		response = handleClaim(cmd, params)
	case "drop":
		response = handleDrop(cmd, params)
	default:
		response = handleUnknown(action)
	}
	w.Write([]byte(response))
}

func handleBlank() string {
	return IDKBlank
}

func handleFind(params []string) string {
	spots, err := store.Find()
	if err != nil {
		return NoSpotsAvailable
	}
	var spotIds []string
	for _, s := range spots {
		spotIds = append(spotIds, s.ID)
	}
	sort.Strings(spotIds)
	return fmt.Sprintf(OpenSpotsTemplate, strings.Join(spotIds, ","))
}

func handleRegister(cmd *slack.SlashCommand, params []string) string {
	var spot store.Spot
	var err error
	if len(params) <= 1 {
		return IDKBlank
	}
	if len(params) == 2 {
		spot, err = store.Register(params[1], cmd.UserName, time.Now())
		if err != nil {
			return fmt.Sprintf(SpotDupeRegistrationErrorTemplate, params[1], spot.RegisteredBy)
		}
	}
	if len(params) > 2 {
		openDate, err := time.Parse(store.SpotDateFormat, params[2])
		if err != nil {
			return fmt.Sprintf(SpotDateFormatRegistrationErrorTemplate, params[2])
		}
		fmt.Println(openDate)
		fmt.Println(time.Now())
		if util.BeforeNow(openDate) {
			return fmt.Sprintf(SpotPastDateRegistrationErrorTemplate, params[2])
		}
		spot, err = store.Register(params[1], cmd.UserName, openDate)
		if err != nil {
			return fmt.Sprintf(SpotDupeRegistrationErrorTemplate, params[1], spot.RegisteredBy)
		}
	}
	return fmt.Sprintf(SpotRegisteredTemplate, spot.ID)
}

func handleClaim(cmd *slack.SlashCommand, params []string) string {
	if len(params) < 2 {
		return IDKBlank
	}
	spot, err := store.Claim(params[1], cmd.UserName)
	if err != nil {
		return fmt.Sprintf(SpotClaimErrorTemplate, params[1])
	}
	return fmt.Sprintf(SpotClaimedTemplate, spot.ID)
}

func handleDrop(cmd *slack.SlashCommand, params []string) string {
	if len(params) < 2 {
		return IDKBlank
	}
	if strings.ToLower(params[1]) == "all" {
		store.DropAllRegistrations(cmd.UserName)
		return fmt.Sprintf("All spots registered by %s", cmd.UserName)
	}
	err := store.DropRegistration(params[1], cmd.UserName)
	if err != nil {
		return fmt.Sprintf(SpotDropRegErrorTemplate, params[1])
	}
	return fmt.Sprintf(SpotDropRegTemplate, params[1])
}

func handleHelp() string {
	return HelpText
}

func handleUnknown(action string) string {
	return fmt.Sprintf(IDKTemplate, action)
}
