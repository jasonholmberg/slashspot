package spot

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jasonholmberg/slashspot/internal/data"
	"github.com/jasonholmberg/slashspot/internal/util"
)

const (
	// NotAvailable - Not available
	NotAvailable = "N/A"
)

var lock sync.Mutex

// NewSpot - A Spot constructor
func NewSpot(ID string, registeredBy string, openDate time.Time) data.Spot {
	now := time.Now()
	if openDate.IsZero() {
		openDate = now
	}
	return data.Spot{
		ID:           ID,
		OpenDate:     openDate.Format(util.SpotDateFormat),
		RegDate:      now.Format(util.SpotDateFormat),
		RegisteredBy: registeredBy,
	}
}

func formatKey(id string, date time.Time) string {
	return fmt.Sprintf("%v-%s", id, date.Format(util.SpotDateFormat))
}

// Find - fins all available spots
func Find() (map[string]data.Spot, error) {
	lock.Lock()
	defer data.Save()
	defer lock.Unlock()
	openSpots := make(map[string]data.Spot)
	store, err := data.Load()
	if err != nil {
		return openSpots, errors.New("error loading spot data")
	}
	log.Println("Finding open spots for today")
	for k, spot := range store {
		if util.BeforeNow(spot.OpenDate) {
			log.Printf(">Cleaning up old registration Id: %v, registered by %v for date: %v", spot.ID, spot.RegisteredBy, spot.OpenDate)
			log.Println()
			data.Delete(k)
			continue
		}
		if util.AfterNow(spot.OpenDate) {
			log.Printf(">Skipping Id: %v, registered by %v for date: %v", spot.ID, spot.RegisteredBy, spot.OpenDate)
			continue
		}
		openSpots[k] = spot
	}
	if len(openSpots) == 0 {
		return openSpots, errors.New("no spots available")
	}
	return openSpots, nil
}

// Claim - claim a spot
func Claim(id string, user string) (data.Spot, error) {
	lock.Lock()
	defer data.Save()
	defer lock.Unlock()
	store, err := data.Load()
	if err != nil {
		return data.Spot{}, errors.New("error loading spot data")
	}
	now := time.Now()
	claimKey := formatKey(id, now)
	for k, spot := range store {
		if k == claimKey {
			log.Printf("Spot %v claimed by %v", id, user)
			data.Delete(k)
			return spot, nil
		}
	}
	return data.Spot{
		ID: NotAvailable,
	}, fmt.Errorf("spot %v not available", id)
}

// Register - register a spot
func Register(id string, user string, openDate time.Time) (data.Spot, error) {
	lock.Lock()
	defer data.Save()
	defer lock.Unlock()
	store, err := data.Load()
	if err != nil {
		return data.Spot{}, errors.New("error loading spot data")
	}
	newSpot := NewSpot(id, user, openDate)
	for k, spot := range store {
		if k == newSpot.Key() {
			return spot, fmt.Errorf("spot %v already registered", id)
		}
	}
	log.Printf("Registered Id: %v by %v for date: %v", newSpot.ID, newSpot.RegisteredBy, newSpot.OpenDate)
	data.Insert(newSpot)
	return newSpot, nil
}

// DropRegistration - drop a registration
func DropRegistration(id string, user string) error {
	lock.Lock()
	defer data.Save()
	defer lock.Unlock()
	store, err := data.Load()
	if err != nil {
		return errors.New("error loading spot data")
	}
	for k, spot := range store {
		if spot.ID == id && spot.RegisteredBy == user {
			data.Delete(k)
			return nil
		}
	}
	return fmt.Errorf("drop reg error of ID: %v", id)
}

// DropAllRegistrations - drop all the registrations for current user
func DropAllRegistrations(user string) {
	lock.Lock()
	defer data.Save()
	defer lock.Unlock()
	store, _ := data.Load()
	for k, spot := range store {
		if spot.RegisteredBy == user {
			data.Delete(k)
		}
	}
}
