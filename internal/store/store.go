package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jasonholmberg/slashspot/internal/util"
)

const (
	// NotAvailable - Not available
	NotAvailable = "N/A"

	// SpotDateFormat - the default spot date formate
	SpotDateFormat = "2006-01-02"
)

var lock sync.Mutex
var spotData *SpotData

type (
	// Spot - a simple spot type
	Spot struct {
		// ID - The spot identifier
		ID string

		// OpenDate - The available date of the spot
		OpenDate string

		// RegDate - The reg date of the spot
		RegDate string

		// RegisteredBy - The user who registered the spot
		RegisteredBy string
	}

	// SpotData - the spot database
	SpotData struct {
		store map[string]Spot
	}
)

// Open - open the spot store
func Open() {
	spotData = &SpotData{
		store: make(map[string]Spot),
	}
	load()
}

// NewSpot - A Spot constructor
func NewSpot(ID string, registeredBy string, openDate time.Time) Spot {
	now := time.Now()
	if openDate.IsZero() {
		openDate = now
	}
	return Spot{
		ID:           ID,
		OpenDate:     openDate.Format(SpotDateFormat),
		RegDate:      now.Format(SpotDateFormat),
		RegisteredBy: registeredBy,
	}
}

func (spot Spot) key() string {
	return fmt.Sprintf("%v-%s", spot.ID, spot.OpenDate)
}

func formatKey(id string, date time.Time) string {
	return fmt.Sprintf("%v-%s", id, date.Format(SpotDateFormat))
}

// IsOpen - data store is open
func IsOpen() bool {
	return spotData.store != nil
}

// Find - fins all available spots
func Find() (map[string]Spot, error) {
	lock.Lock()
	defer save()
	defer lock.Unlock()
	load()
	openSpots := make(map[string]Spot)
	for k, spot := range spotData.store {
		spotDate, _ := time.Parse(SpotDateFormat, spot.OpenDate)
		if util.BeforeNow(spotDate) {
			log.Printf("Cleaning up old registration Id: %v, registered by %v for date: %v", spot.ID, spot.RegisteredBy, spot.OpenDate)
			delete(spotData.store, k)
			continue
		}
		if util.AfterNow(spotDate) {
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
func Claim(id string, user string) (Spot, error) {
	lock.Lock()
	defer save()
	defer lock.Unlock()
	now := time.Now()
	claimKey := formatKey(id, now)
	for k, spot := range spotData.store {
		if k == claimKey {
			log.Printf("Spot %v claimed by %v", id, user)
			delete(spotData.store, k)
			return spot, nil
		}
	}
	return Spot{
		ID: NotAvailable,
	}, fmt.Errorf("spot %v not available", id)
}

// Register - register a spot
func Register(id string, user string, openDate time.Time) (Spot, error) {
	lock.Lock()
	defer save()
	defer lock.Unlock()
	load()
	newSpot := NewSpot(id, user, openDate)
	for k, spot := range spotData.store {
		if k == newSpot.key() {
			return spot, fmt.Errorf("spot %v already registered", id)
		}
	}
	spotData.store[newSpot.key()] = newSpot
	return newSpot, nil
}

// DropRegistration - drop a registration
func DropRegistration(id string, user string) error {
	lock.Lock()
	defer save()
	defer lock.Unlock()
	load()
	for k, spot := range spotData.store {
		if spot.ID == id && spot.RegisteredBy == user {
			delete(spotData.store, k)
			return nil
		}
	}
	return fmt.Errorf("drop reg error of ID: %v", id)
}

// DropAllRegistrations - drop all the registrations for current user
func DropAllRegistrations(user string) {
	lock.Lock()
	defer save()
	defer lock.Unlock()
	load()
	for k, spot := range spotData.store {
		if spot.RegisteredBy == user {
			delete(spotData.store, k)
		}
	}
}

// saves the spot stroe
func save() error {
	// Save to disk
	os.MkdirAll(os.Getenv("SPOT_DATA_DIR"), os.ModePerm)
	f, err := os.Create(DataFilePath())
	defer f.Close()
	if err != nil {
		log.Printf("Error saving spot store %v", err)
	}
	if spotData.store == nil {
		spotData.store = make(map[string]Spot)
		return nil
	}
	r, err := marshal(spotData.store)
	if err != nil {
		errMsg := "Error marshalling spot-store"
		return errors.New(errMsg)
	}
	_, err = io.Copy(f, r)
	return err
}

func load() error {
	// Load from disk
	f, err := os.Open(DataFilePath())
	defer f.Close()
	if err != nil {
		log.Print("No data file to load, creating one")
		save()
	}
	unmarshal(f, &spotData.store)
	return nil
}

// DataFilePath - path to data file
func DataFilePath() string {
	return filepath.Join(os.Getenv("SPOT_DATA_DIR"), os.Getenv("SPOT_DATA_FILE"))
}

func marshal(spots map[string]Spot) (io.Reader, error) {
	b, err := json.MarshalIndent(spots, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func unmarshal(r io.Reader, data *map[string]Spot) error {
	return json.NewDecoder(r).Decode(data)
}
