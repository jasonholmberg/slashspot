package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const (
	add = "add"
	drop = "drop"
)

var store map[string]Spot
var lock sync.Mutex

// Open - open the spot store
func Open() {
	store = make(map[string]Spot)
	Load()
}

// IsOpen - data store is open
func IsOpen() bool {
	return store != nil
}

func save() error {
	// Save to disk
	os.MkdirAll(os.Getenv("SPOT_DATA_DIR"), os.ModePerm)
	f, err := os.Create(FilePath())
	defer f.Close()
	if err != nil {
		log.Printf("Error saving spot store %v", err)
	}
	if store == nil {
		store = make(map[string]Spot)
		return nil
	}
	r, err := marshal(store)
	if err != nil {
		errMsg := "Error marshalling spot-store"
		return errors.New(errMsg)
	}
	_, err = io.Copy(f, r)
	return err
}

// Load - load the data file
func Load() (map[string]Spot, error) {
	// Load from disk
	f, err := os.Open(FilePath())
	defer f.Close()
	if err != nil {
		log.Print("No data file to load, creating one")
		save()
	}
	unmarshal(f, &store)
	return store, nil
}

// Drop - drop the spot
func Drop(s Spot) {
	persist(s, drop)
}

// Add - add the spot
func Add(s Spot) {
	persist(s, add)
}

// Persist - applys changes to the map and saves
func persist(s Spot, op string) {
	lock.Lock()
	defer lock.Unlock()
	defer save()
	if op == add {
		store[s.Key()] = s
	}
	if op == drop {
		delete(store, s.Key())
	}
}

// FilePath - path to data file
func FilePath() string {
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
